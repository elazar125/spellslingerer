class SsrrList extends HTMLElement {
  #page = 1
  #keepLoading = true
  #isLoading = false
  #list = null
  #filterFields = null
  #sortFields = null
  #addFilterSelect = null
  #initialListContent = ''
  #recordType = ''
  #options = {}
  #defaultFilter = null
	#filterHtml = {}
  #sortOptions ={}

  constructor() {
    super()

    this.#initDom()
    this.#initialListContent = this.#list.innerHTML
    if (this.#recordType !== '') {
      this.initLoader()
    }
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

  set filterHtml(value) {
    this.#filterHtml = value
    this.#addFilterSelect.innerHTML = Object.keys(this.#filterHtml)
      .map((key) => `<option value="${key}">${key}</option>`).join()
  }

  set sortOptions(value) {
    this.#sortOptions = value
  }

  renderer(record) { }
  postRenderer(record) { }

  #initDom() {
    const shadow = this.attachShadow({ mode: 'open' })
		const wrapper = document.createElement('section')
		const details = document.createElement('details')
		const summary = document.createElement('summary')
		summary.innerText = 'Filter and Sort Options'
		const form = document.createElement('form')

    const filterWrapper = document.createElement('div')
		this.#filterFields = document.createElement('div')
		this.#filterFields.classList.add('grid')
		this.#filterFields.id = 'filter-fields'
		this.#filterFields.style = '--grid-min-width: 25ch'
    const filterTitleWrapper = document.createElement('div')
    filterTitleWrapper.classList.add('flex-apart')
		const filterTitle = document.createElement('h2')
		filterTitle.innerText = 'Filter'
    const addFilterForm = document.createElement('form')
    addFilterForm.classList.add('flex')
    this.#addFilterSelect = document.createElement('select')
    const addFilterButton = document.createElement('button')
    addFilterButton.classList.add('small-button')
    addFilterButton.innerText = '+'
    addFilterButton.type = 'button'
    addFilterButton.addEventListener('click', () => {
      this.addFilter(this.#addFilterSelect.value)
    })
    addFilterForm.appendChild(this.#addFilterSelect)
    addFilterForm.appendChild(addFilterButton)
    filterTitleWrapper.appendChild(filterTitle)
    filterTitleWrapper.appendChild(addFilterForm)
		filterWrapper.appendChild(filterTitleWrapper)
    filterWrapper.appendChild(this.#filterFields)

    const sortWrapper = document.createElement('div')
		this.#sortFields = document.createElement('div')
		this.#sortFields.classList.add('grid')
		this.#sortFields.id = 'sort-fields'
		this.#sortFields.style = '--grid-min-width: 25ch'
    const sortTitleWrapper = document.createElement('div')
    sortTitleWrapper.classList.add('flex-apart')
		const sortTitle = document.createElement('h2')
		sortTitle.innerText = 'Sort'
    const addSortForm = document.createElement('form')
    addSortForm.classList.add('flex')
    const addSortButton = document.createElement('button')
    addSortButton.classList.add('small-button')
    addSortButton.innerText = '+'
    addSortButton.type = 'button'
    addSortButton.addEventListener('click', this.addSort.bind(this))
    addSortForm.appendChild(addSortButton)
    sortTitleWrapper.appendChild(sortTitle)
    sortTitleWrapper.appendChild(addSortForm)
		sortWrapper.appendChild(sortTitleWrapper)
    sortWrapper.appendChild(this.#sortFields)

		const formControls = document.createElement('div')
		formControls.classList.add('flex')
		const submitButton = document.createElement('button')
		submitButton.type = 'submit'
		submitButton.innerText = 'Search'
		const resetButton = document.createElement('button')
		resetButton.type = 'reset'
		resetButton.innerText = 'Reset'
		formControls.appendChild(submitButton)
		formControls.appendChild(resetButton)

    this.#list = document.createElement('section')
    this.#list.classList.add('grid')
		this.#list.style = '--grid-min-width: 150px'

    const style = document.createElement('style')
    style.innerHTML = this.#getStyle()

		form.appendChild(filterWrapper)
		form.appendChild(sortWrapper)
		form.appendChild(formControls)
		details.appendChild(summary)
		details.appendChild(form)
		wrapper.appendChild(details)
    wrapper.appendChild(style)
		shadow.appendChild(wrapper)
    shadow.appendChild(this.#list)

    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      this.updateFilter()
    })
  }

	addFilter(name) {
    const filterWrapper = document.createElement('div')
    filterWrapper.classList.add('input-wrapper')
    const removeButton = document.createElement('button')
    removeButton.type = 'button'
    removeButton.innerText = 'x'
    removeButton.classList.add('small-button', 'remove-button')
    removeButton.addEventListener('click', () => this.#filterFields.removeChild(filterWrapper))
    filterWrapper.innerHTML = this.#filterHtml[name]
    filterWrapper.appendChild(removeButton)
		this.#filterFields.appendChild(filterWrapper)
	}

	addSort(initialValue, initialDir){
    const count = this.#sortFields.childElementCount
    const sortWrapper = document.createElement('div')
    sortWrapper.classList.add('input-group', 'input-wrapper')
    const removeButton = document.createElement('button')
    removeButton.type = 'button'
    removeButton.innerText = 'x'
    removeButton.classList.add('small-button', 'remove-button')
    removeButton.addEventListener('click', () => this.#sortFields.removeChild(sortWrapper))
    sortWrapper.innerHTML = `
	    <label for="prop${count}">Property</label>
	    <select id="prop${count}" data-main-input>
	      <option value="">Unset</option>
        ${
          Object.entries(this.#sortOptions)
            .map(([name, value]) => {
              return `<option value="${value}" ${name === initialValue ? 'selected' : ''}>${name}</option>`
            })
            .join()
        }
	    </select>
	    <div>
	      <label for="dir${count}">Direction</label>
	      <select id="dir${count}" data-opcode>
	        <option value="+" ${initialDir === 'asc' ? 'selected' : ''}>Ascending</option>
	        <option value="-" ${initialDir === 'desc' ? 'selected' : ''}>Descending</option>
	      </select>
	    </div>
  	`
    sortWrapper.appendChild(removeButton)
    this.#sortFields.appendChild(sortWrapper)
  }

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

  async initLoader() {
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

  #getStyle() {
    return `
 :host {
  font-size: 1rem;
  font-family: sans-serif;
  line-height: 1.4;
  background-color: #1c1b22;
  color: #fbfbfe;

  --gutter: clamp(1rem, 5vw, 2rem);
  --menu-size: 185px;
  --header-size: 60px;

  display: grid;
  gap: var(--gutter);

  --primary-colour: rgb(154, 81, 89);
  --primary-colour-highlighted: rgb(191, 120, 128);
  --secondary-colour: rgb(255, 132, 0);
  --secondary-colour-dimmed: rgb(79, 57, 33);
  --unaccented-colour: rgb(69, 69, 69);
  --unaccented-colour-highlighted: rgb(99, 99, 99);

  --control-bg-colour: rgb(54, 54, 54);
  --control-bg-highlight: rgb(70, 70, 70);
  --control-border-radius: 5px;
  --control-padding: 1.5rem;
  --control-shadow:
    0px 2px 4px -1px rgb(0 0 0 / 0.2),
    0px 4px 5px 0px rgb(0 0 0 / 0.14),
    0px 1px 10px 0px rgb(0 0 0 / 0.12);
  --text-shadow: 0.1em 0.1em 0.2em rgb(0 0 0 / 0.4);

  --ssrr-deck-padding: var(--control-padding);
  --ssrr-deck-background-colour: var(--control-bg-colour);
  --ssrr-deck-border-radius: var(--control-border-radius);
  --ssrr-deck-box-shadow: var(--control-shadow);
  --ssrr-deck-card-colour-white: rgb(164, 154, 113);
  --ssrr-deck-card-colour-blue: rgb(26, 26, 152);
  --ssrr-deck-card-colour-black: rgb(26, 26, 26);
  --ssrr-deck-card-colour-red: rgb(153, 26, 26);
  --ssrr-deck-card-colour-green: rgb(26, 98, 26);
  --ssrr-deck-card-colour-colourless: rgb(109, 109, 109);

  --colour-error: hsl(4, 90%, 48%);
  --colour-info: hsl(175, 40%, 20%);
  --colour-success: hsl(106, 39%, 27%);
  --colour-warning: hsl(39, 86%, 36%);
}

section {
  background-color: var(--control-bg-colour);
  padding: var(--control-padding);
  border-radius: var(--control-border-radius);
  box-shadow: var(--control-shadow);
  height: max-content;
  word-wrap: break-word;
  min-width: 0;
}

summary {
  cursor: pointer;
}

.grid {
  display: grid;
  gap: var(--gutter);
  grid-template-columns: repeat(auto-fit, minmax(var(--grid-min-width, 200px), 1fr));
}

.flex {
  display: flex;
  gap: var(--gutter);
}

.flex-apart {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.tile-thumbnail {
  max-width: 4rem;
}

.tile {
  background-image: var(--deck-tile-image);
  border-radius: var(--control-border-radius);
  box-shadow: var(--control-shadow);
  aspect-ratio: 4/3;
  width: 200px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  cursor: pointer;
  background-size: cover;
  position: relative;
  isolation: isolate;
  overflow: hidden;
}

.tile::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  inset: 0;
  background-image: linear-gradient(to bottom, #000000, transparent, #000000);
  z-index: -1;
}

.tile__details {
  height: 2rem;
  display: flex;
  justify-content: space-between;
}

.tile__icons {
  height: 3rem;
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.tile svg {
  max-width: 2rem;
  max-height: 1rem;
  fill: currentColor;
}

button:hover, .button:hover {
  background-color: var(--primary-colour-highlighted);
}

button:active, .button:active {
  background-color: var(--primary-colour-highlighted);
}

button:focus-visible, .button:focus-visible {
  background-color: var(--primary-colour-highlighted);
}

button, .button {
  text-decoration: unset;
  color: unset;
  border: none;
  width: max-content;
  padding: 1rem 2rem;
  background-color: var(--button-colour, var(--primary-colour));
  border-radius: var(--control-border-radius);
  text-transform: uppercase;
  box-shadow: var(--control-shadow);
  font-weight: bold;
  font-size: inherit;
  font-family: inherit;
  line-height: normal;
  letter-spacing: 2px;
  cursor: pointer;
}

.input-wrapper {
  position: relative;
}

button.small-button {
  padding: 3px 9px;
}

button.remove-button {
  position: absolute;
  top: 0;
  right: 0;
  background-color: var(--unaccented-colour);
}

button.remove-button:active {
  background-color: var(--unaccented-colour-highlighted);
}

button.remove-button:hover {
  background-color: var(--unaccented-colour-highlighted);
}

form {
  display: grid;
  gap: var(--gutter);
  padding: var(--control-padding);
}

@media(min-width: 800px) {
  .inline-form {
    grid-auto-flow: column;
    grid-template-columns: repeat(3, 1fr);
  }

  .inline-form button {
    width: 100%;
    max-width: 230px;
  }

  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
}
@media(max-width: 800px) {
  .input-group {
    max-width: 30ch;
  }
}

input[type="text"],
input[type="number"],
select:not([data-opcode]),
fieldset {
  width: 100%;
  min-width: 0;
}

input {
  accent-color: var(--primary-colour);
  font-size: 1rem;
}

textarea {
  accent-color: var(--primary-colour);
  font-size: 1rem;
}

fieldset {
  border: none;
  margin: 0;
  padding: 0;
}

/* TODO: Focus and active styles */

.card {
  display: flex;
  flex-direction: column;
  place-items: center;
}

.card__buttons {
  display: flex;
  place-items: center;
  gap: 1rem;
}

.card button {
  padding: 12px;
  border-radius: 50%;
  aspect-ratio: 1;
}

button[aria-label="remove card"] {
  padding: 15px;
}

img {
  max-width: 100%;
}

textarea {
  width: 100%;
  height: 100%;
}
      `
  }
}

export default SsrrList
