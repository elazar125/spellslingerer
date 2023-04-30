import { confirmModal, alertOnErrorAsync } from '/js/main.js'

class Deck {
  constructor(
    originalDeck,
    ssrrDeck,
    manaCostChart,
    rarityChart,
    cardTypeChart,
    deckTitle,
    deckDescription,
    isPublic,
    currentCount,
  ) {
    this.deck = JSON.parse(JSON.stringify(originalDeck))
    this.originalDeck = JSON.parse(JSON.stringify(originalDeck))
    if (!this.deck.expand) this.deck.expand = {}
    if (!this.deck.card_details) this.deck.card_details = {}
    if (!this.deck.cards) this.deck.cards = []
    if (!this.deck.expand.cards) this.deck.expand.cards = []
    this.ssrrDeck = ssrrDeck
    this.manaCostChart = manaCostChart
    this.rarityChart = rarityChart
    this.cardTypeChart = cardTypeChart
    this.deckTitle = deckTitle
    this.deckDescription = deckDescription
    this.isPublic = isPublic
    this.currentCount = currentCount
    this.#render()
  }

  #render() {
    this.ssrrDeck.render(this.deck)
    this.manaCostChart.render({ type: 'bar', title: 'Mana Cost', data: this.#formatManaCostData(this.deck) })
    this.rarityChart.render({ type: 'bar', title: 'Rarity', data: this.#formatRarityData(this.deck) })
    this.cardTypeChart.render({ type: 'pie', title: 'Card Type', data: this.#formatCardTypeData(this.deck) })
    this.deckTitle.value = this.deck.name
    this.deckDescription.value = this.deck.description
    this.isPublic.checked = this.deck.is_public
  }

  setName(name) {
    this.deck.name = name
    this.#render()
  }

  setDescription(description) {
    this.deck.description = description
    this.#render()
  }

  setImage(image) {
    this.deck.image = image
    this.#render()
  }

  setIsPublic(isPublic) {
    this.deck.is_public = isPublic
    this.#render()
  }

  hasSpellslinger() {
    return !!this.deck.spellslinger
  }
 
  canSetSplash() {
    if (!this.deck.spellslinger) return false

    return [
      'Gideon',
      'Jace',
      'Liliana',
      'Chandra',
      'Sorin',
      'Yanling',
    ].includes(this.deck.expand.spellslinger.name)
  }
 
  setSplash(data) {
    if (data === '') return
    const splash = JSON.parse(data)
  
    this.deck.splash = splash.id
    this.deck.expand.splash = splash
    this.#render()
  }

  setSpellslinger(data) {
    const spellslinger = JSON.parse(data)

    this.deck.spellslinger = spellslinger.id
    this.deck.expand.spellslinger = spellslinger

    this.#render()
  }

  #updateCurrentCardCount(toAdd) {
    let count = parseInt(this.currentCount.innerHTML)
    count += toAdd
    this.currentCount.innerHTML = count
  }

  removeCard(data) {
    const card = JSON.parse(decodeURIComponent(data))

    if (this.deck.card_details[card.id]) {
      const quantity = this.deck.card_details[card.id].standard_quantity

      if (quantity > 1) {
        this.deck.card_details[card.id].standard_quantity--
      }
      else {
        delete this.deck.card_details[card.id]
        this.deck.cards = this.deck.cards.filter((c) => c != card.id)
        this.deck.expand.cards = this.deck.expand.cards.filter((c) => c.id != card.id)
      }
      this.#updateCurrentCardCount(-1)
    }

    this.#render()
  }

  addCard(data) {
    const card = JSON.parse(decodeURIComponent(data))

    if (card.expand.type.name.toLowerCase() === 'land') {
      this.deck.land = card.id
      this.deck.expand.land = card
    }
    else if (this.deck.card_details[card.id]) {
      if (this.deck.card_details[card.id].standard_quantity < (card.legendary ? 1 : 2)) {
        this.deck.card_details[card.id].standard_quantity++

        this.#updateCurrentCardCount(1)
      }
    }
    else {
      this.deck.card_details[card.id] = { name: card.name, foil_quantity: 0, standard_quantity: 1 }
      this.deck.cards.push(card.id)
      this.deck.expand.cards.push(card)

      this.#updateCurrentCardCount(1)
    }

    this.#render()
  }
  
  getCards() {
    return this.deck.expand.cards
  }

  reset() {
    confirmModal({
      titleText: 'Revert unsaved changes',
      bodyText: 'You are about to reset all unsaved work. Are you sure you want to do so?',
      onConfirm: () => {
        this.deck = JSON.parse(JSON.stringify(this.originalDeck))
        this.#render()
      },
      onCancel: () => { },
    })
  }
  
  save() {
    alertOnErrorAsync(async () => {
      if (this.deck.id) {
        const record = await client.collection('decks').update(this.deck.id, this.deck)
        location.href = `/decks/${record.id}`
      }
      else {
        const record = await client.collection('decks').create(this.deck)
        location.href = `/decks/${record.id}`
      }
    })
  }

  delete() {
    if (this.deck.id) {
      confirmModal({
        titleText: 'Are you sure?',
        bodyText: 'You are about to delete this deck. This can not be undone. Are you sure?',
        onConfirm: () => alertOnErrorAsync(async () => {
          await client.collection('decks').delete(this.deck.id)
          location.href = '/my-decks/'
        }),
        onCancel: () => { },
      })
    }
    else {
      confirmModal({
        titleText: 'Leaving unsaved page',
        bodyText: 'You are about to be redirected to the "My Decks" page and lose all unsaved work. Are you sure you want to leave?',
        onConfirm: () => {
          location.href = '/my-decks/'
        },
        onCancel: () => { },
      })
    }
  }
  
  #formatManaCostData(deck) {
    const result = []

    this.deck.expand.cards.forEach((c) => {
      if (!result[c.cost]) result[c.cost] = { label: `${c.cost}`, value: 0, colour: 'rgb(153, 80, 88)' }
      result[c.cost].value += this.#getQuantity(c.id)
    })

    result.forEach((r, i) => { if (!r) result[i] = { label: `${i}`, value: 0, colour: 'rgb(153, 80, 88)' } })

    return result
  }

  #formatRarityData(deck) {
    const result = {
      core: { label: 'Core', value: 0, colour: 'rgb(150,150,150)' },
      signature: { label: 'Signature', value: 0, colour: 'rgb(187, 172, 41)' },
      common: { label: 'Common', value: 0, colour: 'rgb(187,187,187)' },
      rare: { label: 'Rare', value: 0, colour: 'rgb(57, 20, 100)' },
      epic: { label: 'Epic', value: 0, colour: 'rgb(136, 11, 115)' },
      mythic: { label: 'Mythic', value: 0, colour: 'rgb(177, 121, 18)' },
    }

    this.deck.expand.cards.forEach((c) => result[c.expand.rarity.name.toLowerCase()].value += this.#getQuantity(c.id))

    return Object.values(result)
  }

  #formatCardTypeData(deck) {
    const result = {
      creature: { label: 'Creature', value: 0, colour: 'rgb(126, 76, 36)' },
      spell: { label: 'Spell', value: 0, colour: 'rgb(42, 119, 108)' },
      trap: { label: 'Trap', value: 0, colour: 'rgb(30, 30, 30)' },
      artifact: { label: 'Artifact', value: 0, colour: 'rgb(124, 76, 68)' },
      skill: { label: 'Skill', value: 0, colour: 'rgb(215, 122, 73)' },
    }

    this.deck.expand.cards.forEach((c) => result[c.expand.type.name.toLowerCase()].value += this.#getQuantity(c.id))

    return Object.values(result)
  }

  #getQuantity(cardId) {
    const card_details = this.deck.card_details[cardId]
    return card_details.standard_quantity + card_details.foil_quantity
  }
}

export default Deck