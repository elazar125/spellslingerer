class Loader {
  #page = 1
  #keepLoading = true
  #isLoading = false
  #list = null
  #filterFields = null
  #sortFields = null
  #initialListContent = ''
  #recordType = ''
  #options = {}
  #defaultFilter = null

  constructor(list, filterFields, sortFields) {
    this.#list = list
    this.#filterFields = filterFields
    this.#sortFields = sortFields
    this.#initialListContent = list.innerHTML
  }

  set recordType(value) {
    this.#recordType = value
  }

  set expandFields(value) {
    this.#options.expand = value
  }

  set defaultFilter(value) {
    this.#defaultFilter = value
  }

  renderer(record) { }
  postRenderer(record) { }

  #loadFilter() { }

  #parseFilter() {
    const filter = this.#parseFilterFields()
    const sort = this.#parseSortFields()

    if (filter) {
      this.#options.filter = filter
    }
    else {
      delete this.#options.filter
    }

    if (sort) {
      this.#options.sort = sort
    }
    else {
      delete this.#options.sort
    }
  }

  async updateFilter() {
    this.#parseFilter()
    this.#clearRecords()
    await this.#loadRecords()
  }

  #clearRecords() {
    this.#page = 1
    this.#list.innerHTML = this.#initialListContent
  }

  async #loadRecords() {
    if (!this.#options?.filter && this.#defaultFilter) {
      this.#options.filter = this.#defaultFilter
    }

    try {
      this.#isLoading = true
      const result = await client.collection(this.#recordType).getList(this.#page, 30, this.#options)
      this.#page++
      this.#keepLoading = result.page < result.totalPages
      this.#list.innerHTML += result.items.map(this.renderer).join('')
      this.postRenderer()
    }
    catch {}
    finally {
      this.#isLoading = false
    }
  }

  async init() {
    this.#loadFilter()
    this.#parseFilter()
    await this.#loadRecords()

    window.onscroll = async () => {
      if (this.#scrolledToBottom() && this.#keepLoading && !this.#isLoading) {
        await this.#loadRecords()
      }
    }
  }

  #scrolledToBottom() {
    const docHeight = Math.max(
      document.body.scrollHeight,
      document.documentElement.scrollHeight,
      document.body.offsetHeight,
      document.documentElement.offsetHeight,
      document.body.clientHeight,
      document.documentElement.clientHeight,
    )
    const scrollTop = Math.ceil(Math.max(
      window.pageYOffset,
      document.body.scrollTop,
      document.documentElement.scrollTop,
    ))
    return innerHeight + scrollTop >= docHeight
  }

  #parseFilterFields() {
    const results = []

    if (this.#defaultFilter) results.push(this.#defaultFilter)

    this.#filterFields.querySelectorAll('.input-group').forEach((group) => {
      const input = group.querySelector('[data-main-input]')
      const opcode = group.querySelector('[data-opcode]')
      const fieldName = group.attributes.getNamedItem('data-field-name')?.value
      const additionalFilterRequired = group.querySelector('[data-additional-filter]')
      const additionalFilter = additionalFilterRequired?.attributes.getNamedItem('data-additional-filter')

      if (input.type === 'fieldset') {
        const fieldsetFilter = this.#parseFieldsetFilter(input, opcode, fieldName)
        if (fieldsetFilter) results.push(fieldsetFilter)
      }
      else if (input.value) {
        switch (input.type) {
          case 'text':
            results.push(...input.value.split(' ').map((s) => `${fieldName}?~'${s}'`))
            break
          case 'number':
            results.push(`${fieldName}${opcode.value}${input.value}`)
            break
          case 'select-one':
            results.push(`${fieldName}=${input.value}`)
            break
          case 'select-multiple':
            results.push(`(${[...input.selectedOptions].map((op) => `${fieldName}='${op.value}'`).join('||')})`)
            break
        }
      }
      if (additionalFilterRequired && additionalFilterRequired.checked && additionalFilter) {
        results.push(additionalFilter.nodeValue)
      }
    })

    return results.join('&&')
  }

  #parseFieldsetFilter(fieldset, opcode, fieldName) {
    const selectedValues = []
    let result = null
    fieldset.querySelectorAll('input').forEach((check) => {
      selectedValues.push([check.value, check.checked])
    })

    if (selectedValues.every(([value, include]) => !include)) {
      return null
    }

    switch(opcode?.value) {
      case 'exact':
        result = '(' + selectedValues
          .filter(([value, include]) => include)
          .map(([value, include]) => `${fieldName}?='${value}'`)
          .join('||') + ')'
        result += '(' + selectedValues
          .filter(([value, include]) => !include)
          .map(([value, include]) => `${fieldName}!~'${value}'`)
          .join('&&') + ')'
        result += `(${fieldName.slice(0, fieldName.lastIndexOf('.'))}:length=${selectedValues.filter(([value, include]) => include).length})`
        break
      case 'include':
        result = selectedValues
          .filter(([value, include]) => include)
          .map(([value, include]) => `${fieldName}?='${value}'`)
          .join('||')
        break
      case 'at-most':
        result = selectedValues
          .filter(([value, include]) => !include)
          .map(([value, include]) => `${fieldName}!~'${value}'`)
          .join('&&')
        break
      default:
        result = selectedValues
          .filter(([value, include]) => include)
          .map(([value, include]) => `${fieldName}?='${value}'`)
          .join('||')
        break
    }

    return result ? `(${result})` : null
  }

  #parseSortFields() {
    const results = []

    this.#sortFields.querySelectorAll('.input-group').forEach((group) => {
      const input = group.querySelector('[data-main-input]')
      const opcode = group.querySelector('[data-opcode]')
      const matchesColour = input.value.match(/(.*colour).sort_order/)

      if (matchesColour) {
        results.push(`${opcode.value}${matchesColour[1]}:length`)
      }

      if (input.value) {
        results.push(`${opcode.value}${input.value}`)
      }
    })

    return results.join(',')
  }
}

export default Loader
