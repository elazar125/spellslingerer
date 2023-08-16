class SsrrDeck extends HTMLElement {
  constructor() {
    super()
    const shadow = this.attachShadow({ mode: 'open' })
    const wrapper = document.createElement('div')
    this.tile = document.createElement('div')
    this.tile.setAttribute('part', 'tile')
    this.cards = document.createElement('div')
    this.cards.classList.add('cards')
    this.cards.setAttribute('part', 'cards')
    const style = document.createElement('style')
    style.innerHTML = this.#getStyle()
    wrapper.appendChild(this.tile)
    wrapper.appendChild(this.cards)

    if(this.hasAttribute('display-code')) {
      this.code = document.createElement('div')
      wrapper.appendChild(this.code)
    }

    wrapper.appendChild(style)
    shadow.appendChild(wrapper)
    this.domain = ''
    
    if(this.hasAttribute('url')) {
      const { domain, deckid } = this.#parseUrl(this.getAttribute('url'))
      this.domain = domain
      this.#loadDeck(deckid)
    }
    if(this.hasAttribute('deck')) {
      this.domain = 'https://spellslingerer.com'
      this.#loadDeck(this.getAttribute('deck'))
    }
  }
  
  #parseUrl(inputUrl) {
    const inputPattern = /^(?<domain>[a-zA-Z0-9:./]+)\/decks\/(?<deckid>[a-zA-Z0-9]{15})$/
    const matchGroups = inputUrl.match(inputPattern)?.groups
    return {
      domain: matchGroups?.domain,
      deckid: matchGroups?.deckid,
    }
  }
  
  async #loadDeck(deckid) {
    const fetchUrl = `${this.domain}/api/collections/decks/records/${deckid}?expand=cards.set,cards.colour,cards.rarity,cards.type,land.set,spellslinger,splash,owner`
    const response = await fetch(fetchUrl)
    const deck = await response.json()
    this.render(deck)
  }

  render(deck) {
    this.tile.innerHTML = this.#renderTile(deck)
    this.cards.innerHTML = this.#renderCardList(deck)
    this.tile.querySelectorAll('#land-image').forEach((card) => {
      const cardHover = this.tile.querySelector('#land-hover')
      card.addEventListener('mouseenter', () => this.#startCardHover(cardHover))
      card.addEventListener('mousemove', (e) => this.#updateCardHover(e, cardHover))
      card.addEventListener('mouseleave', () => this.#endCardHover(cardHover))
    })
    this.cards.querySelectorAll('.card').forEach((card) => {
      const cardHover = card.nextElementSibling
      card.addEventListener('mouseenter', () => this.#startCardHover(cardHover))
      card.addEventListener('mousemove', (e) => this.#updateCardHover(e, cardHover))
      card.addEventListener('mouseleave', () => this.#endCardHover(cardHover))
    })
    if(this.hasAttribute('display-code')) {
      this.code.innerHTML = this.#renderCode(deck)
      this.code.querySelector('button').addEventListener('click', () => {
        const code = this.code.querySelector('p').textContent
        navigator.clipboard.writeText(code)
      })
    }
  }
  
  #renderTile(deck) {
    return `
      <div class="tile" style="--tile-image: url('${this.domain}${deck.image.replace("'", "\\'")}')">
        <h3><a href="${this.domain}/decks/${deck.id || ''}">${deck.name}</a></h3>
        <div class="icons">
          ${deck.splash ? `<img src="${this.domain}/images/mana/${deck.expand.splash.name}.svg" />` : ''}
          <a href="${this.domain}/spellslingers/${deck.expand.spellslinger.id}">
            <img id="spellslinger-image" src="${this.domain}/images/spellslingers/${deck.expand.spellslinger.name}.webp" alt="Spellslinger: ${deck.expand.spellslinger.name}"/>
          </a>
          <a href="${this.domain}/cards/${deck.expand.land.id}">
            <img id="land-image" src="${this.domain}${this.#formatCardImagePath('full_art', deck.expand.land)}" alt="Land: ${deck.expand.land.name}"/>
          </a>
        </div>
      </div>
      <img id="land-hover" class="card__hover" src="${this.domain}${this.#formatCardImagePath('text', deck.expand.land)}" />
    `
  }
  
  #renderCardList(deck) {
    return deck.expand.cards.sort((a, b) => a.cost - b.cost).map((card) => {
      return `
        ${deck.card_details[card.id].standard_quantity > 0 ? `
          <div class="card" style="${this.#generateBackgroundGradient(card.expand.colour)}">
            <span class="card__name"><a href="${this.domain}/cards/${card.id}">${card.name}</a></span>
            <span class="card__details">
              <span class="card__image" style="--card-background-image: url('${this.domain}${this.#formatCardImagePath('tiles', card)}')"></span>
              <span class="card__quantity">x${deck.card_details[card.id].standard_quantity}</span>
            </span>
          </div>
          <img class="card__hover" src="${this.domain}${this.#formatCardImagePath('text', card)}" />
        ` : ''}
        ${deck.card_details[card.id].foil_quantity > 0 ? `
          <div class="card foil" style="${this.#generateBackgroundGradient(card.expand.colour)}">
            <span class="card__name"><a href="${this.domain}/cards/${card.id}">${card.name}</a></span>
            <span class="card__details">
              <span class="card__image" style="--card-background-image: url('${this.domain}${this.#formatCardImagePath('tiles', card)}')"></span>
              <span class="card__quantity">x${deck.card_details[card.id].foil_quantity}</span>
            </span>
          </div>
          <img class="card__hover" src="${this.domain}${this.#formatCardImagePath('foil_text', card)}" />
        ` : ''}
      `
    }).join('')
  }
  
  #renderCode(deck) {
    return `
      <section id="code">
        <div id="code-header" class="flex flex--row">
          <h3>Import Code</h3>
          <button id="code-copy-button">Copy</button>
        </div>
        <p>${deck.code}</p>
      </section>
    `
  }

  #getStyle() {
    return `
      :host {
        --font-colour: var(--ssrr-deck-font-colour, #fbfbfe);
        --text-shadow: var(--ssrr-deck-text-shadow, 0 0 5px black);
        --padding: var(--ssrr-deck-padding, 1.5rem);
        --background-colour: var(--ssrr-deck-background-colour, rgb(54, 54, 54));
        --box-shadow: var(--ssrr-deck-box-shadow);
        --border-radius: var(--ssrr-deck-border-radius, 5px);
        --card-colour-white: var(--ssrr-deck-card-colour-white, rgb(164, 154, 113));
        --card-colour-blue: var(--ssrr-deck-card-colour-blue, rgb(26, 26, 152));
        --card-colour-black: var(--ssrr-deck-card-colour-black, rgb(26, 26, 26));
        --card-colour-red: var(--ssrr-deck-card-colour-red, rgb(153, 26, 26));
        --card-colour-green: var(--ssrr-deck-card-colour-green, rgb(26, 98, 26));
        --card-colour-colourless: var(--ssrr-deck-card-colour-colourless, rgb(109, 109, 109));
        --unaccented-colour: rgb(69, 69, 69);
        --unaccented-colour-highlighted: rgb(99, 99, 99);

        display: block;
        color: var(--font-colour);
        border-radius: var(--border-radius);
        background-color: var(--background-colour);
        box-shadow: var(--box-shadow);
        width: min-content;
        overflow: hidden;
        height: max-content;
      }

      h3 {
        margin: 0;
      }

      .tile {
        background-image: var(--tile-image);
        background-size: cover;
        background-position: center;
        min-width: 300px;
        height: 125px;
        padding: var(--padding);
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        position: relative;
        isolation: isolate;
      }

      .tile::after {
        content: ' ';
        inset: 0;
        background-image: linear-gradient(to bottom, rgb(0 0 0 / 0), var(--background-colour));
        position: absolute;
        z-index: -1;
      }

      .icons {
        height: 3rem;
        display: flex;
        gap: 1rem;
        justify-content: flex-end;
        width: 100%;
      }

      #spellslinger-image {
        height: 3rem;
      }

      #land-image {
        height: 4rem;
        margin-inline: -10px;
        margin-block: -5px;
      }

      .cards {
        padding-block: 1em;
      }

      .card {
        display: block;
        padding: 0.5em var(--padding);
        display: flex;
        justify-content: space-between;
        position: relative;
        margin-block: 3px;
        isolation: isolate;
        background-image: var(--card-colours);
      }

      .card::after {
        content: ' ';
        position: absolute;
        inset: 0;
        background-image: linear-gradient(to right, var(--background-colour), rgb(0 0 0 / 0) 7% 93%, var(--background-colour)), linear-gradient(to bottom, rgb(0 0 0 / 0), rgb(0 0 0 / 0.5));
        z-index: -1;
      }

      .card__details {
        position: relative;
        margin-block: -0.5em;
        display: flex;
        align-items: center;
      }

      .card__details::after {
        content: ' ';
        position: absolute;
        inset: 0;
        background-image: linear-gradient(to bottom, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0.5)), var(--card-colours);
        z-index: 0;
        mask: linear-gradient(to right, rgb(0 0 0 / 1), rgb(0 0 0 / 0) 15%);
      }

      h3 a,
      .card__name a {
        text-decoration: none;
        color: inherit;
        text-shadow: var(--text-shadow);
      }

      .card__image {
        background-image: var(--card-background-image);
        width: 100px;
        height: 100%;
        display: inline-block;
        background-position: center;
        background-size: cover;
      }

      .card__quantity {
        padding-inline-start: 1em;
      }
      
      .card__hover {
        display: none;
        opactiy: 0;
        position: absolute;
        height: 300px;
        z-index: 2;
        top: var(--y, 0);
        left: var(--x, 0);
      }
      
      .card__hover.show {
        display: block;
        opacity: 1;
        transition: opacity 100ms ease;
      }
      
      #code {
        word-break: break-all;
        padding: var(--ssrr-deck-padding);
      }

      #code-header {
        display: flex;
        gap: var(--gutter);
        justify-content: space-between;
        align-items: baseline;
      }

      button, .button {
        text-decoration: unset;
        color: unset;
        border: none;
        width: max-content;
        padding: 1rem 2rem;
        background-color: var(--unaccented-colour);
        border-radius: var(--border-radius);
        text-transform: uppercase;
        box-shadow: var(--box-shadow);
        font-weight: bold;
        letter-spacing: 2px;
        cursor: pointer;
      }
      
      #code-copy-button:hover {
        background-color: var(--unaccented-colour-highlighted);
      }

      #code-copy-button:active {
        background-color: var(--unaccented-colour-highlighted);
      }

      #code-copy-button:focus-visible {
        background-color: var(--unaccented-colour-highlighted);
      }
    `
  }

  #formatCardImagePath(type, card) {
    return `/images/cards/${type}/${card.name}.${ type === 'tiles' ? 'jpeg' : 'webp'}`.replaceAll(' ', '%20').replaceAll('\'', '%27')
  }

  #generateBackgroundGradient(colours) {
    const coloursAsVars = colours.map((c, i, ary) => `var(--card-colour-${c.name.toLowerCase()}) ${i*100/ary.length}% ${(i+1)*100/ary.length}%`);
    return `--card-colours: linear-gradient(to bottom,${coloursAsVars.join(',')})`
  }
  
  #startCardHover(cardHover) {
    cardHover.classList.add('show')
  }
  
  #updateCardHover(e, cardHover) {
    cardHover.style.setProperty('--x', `${e.pageX+10}px`)
    cardHover.style.setProperty('--y', `${e.pageY-150}px`)
  }
  
  #endCardHover(cardHover) {
    cardHover.classList.remove('show')
  }
}

customElements.define('ssrr-deck', SsrrDeck)

export default SsrrDeck
