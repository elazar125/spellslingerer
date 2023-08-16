import SsrrList from '/js/ssrr-list.js'

class SsrrCardList extends SsrrList {
  #addCardCallback = null
  #removeCardCallback = null
  
  constructor() {
    super()
    this.recordType = 'cards'
    this.expandFields = 'set,type,rarity,colour'
    this.filterHtml = this.#filterHtml
    this.sortOptions = this.#sortOptions
    this.#addDefaultFilters()
    this.initLoader()
  }

  set addCardCallback(value) {
    this.#addCardCallback = value
    this.postRenderer()
  }

  set removeCardCallback(value) {
    this.#removeCardCallback = value
    this.postRenderer()
  }

  #addDefaultFilters() {
    this.addFilter('name')
    this.addFilter('colour')
    this.addFilter('rarity')
    this.addSort('Set', 'desc')
    this.addSort('Colour', 'asc')
    this.addSort('Cost', 'asc')
  }

  postRenderer() {
    this.shadowRoot.querySelectorAll('[aria-label="remove card"]').forEach((button) => {
      button.addEventListener('click', (e) => this.#removeCardCallback(e.target.getAttribute('data-card')))
    })
    this.shadowRoot.querySelectorAll('[aria-label="add card"]').forEach((button) => {
      button.addEventListener('click', (e) => this.#addCardCallback(e.target.getAttribute('data-card')))
    })
  }

  renderer(card) {
    return `
    <div class="card">
      <picture>
        <source
          srcset="/images/cards/text/${card.name.replaceAll(' ', '%20')}.webp"
          type="image/webp"
        />
        <source
          srcset="/images/cards/text/${card.name.replaceAll(' ', '%20')}.png"
          type="image/png"
        />
        <img
          loading="lazy"
          srcset="/images/cards/text/${card.name.replaceAll(' ', '%20')}.webp"
          alt="${card.name}"
        />
      </picture>
      <span class="card__buttons">
        <button aria-label="remove card" data-card="${encodeURIComponent(JSON.stringify(card))}">-</button>
        <button aria-label="add card" data-card="${encodeURIComponent(JSON.stringify(card))}">+</button>
      </span>
    </div>
    `
  }

  #sortOptions = {
    'Name': 'name',
    'Ability': 'ability',
    'Cost': 'cost',
    'Power': 'power',
    'Health': 'health',
    'Colour': 'colour.sort_order',
    'Set': 'set.sort_order',
    'Rarity': 'rarity.sort_order',
    'Artist': 'artist',
    'Type': 'type.name',
    'Subtype': 'subtype',
    'Chance': 'chance',
    'Legendary': 'legendary',
  }

	#filterHtml = {
	  'name': `<div class="input-group" data-field-name="name">
	    <label for="name">Name</label>
	    <input type="text" id="name" data-main-input />
	  </div>`,
	  'ability': `<div class="input-group" data-field-name="ability">
	    <label for="ability">Ability</label>
	    <input type="text" id="ability" data-main-input />
	  </div>`,
	  'cost': `<div class="input-group" data-field-name="cost">
	    <label for="cost">Cost</label>
	    <input type="number" id="cost" data-main-input />
	    <div>
	      <label for="cost-op">Operation</label>
	      <select id="cost-op" data-opcode>
	        <option value="=">=</option>
	        <option value="!=">!=</option>
	        <option value=">">&gt;</option>
	        <option value=">=">&gt;=</option>
	        <option value="<">&lt;</option>
	        <option value="<=">&lt;=</option>
	      </select>
	    </div>
	  </div>`,
	  'power': `<div class="input-group" data-field-name="power">
	    <label for="power">Power</label>
	    <input type="number" id="power" data-main-input />
	    <div>
	      <label for="power-op">Operation</label>
	      <select id="power-op" data-opcode>
	        <option value="=">=</option>
	        <option value="!=">!=</option>
	        <option value=">">&gt;</option>
	        <option value=">=">&gt;=</option>
	        <option value="<">&lt;</option>
	        <option value="<=">&lt;=</option>
	      </select>
	    </div>
	  </div>`,
	  'health': `<div class="input-group" data-field-name="health">
	    <label for="health">Health</label>
	    <input type="number" id="health" data-main-input />
	    <div>
	      <label for="health-op">Operation</label>
	      <select id="health-op" data-opcode>
	        <option value="=">=</option>
	        <option value="!=">!=</option>
	        <option value=">">&gt;</option>
	        <option value=">=">&gt;=</option>
	        <option value="<">&lt;</option>
	        <option value="<=">&lt;=</option>
	      </select>
	    </div>
	  </div>`,
	  'colour': `<div class="input-group" data-field-name="colour.name">
	    <fieldset data-main-input>
	      <legend>Colour</legend>
	      <input type="checkbox" id="white" value="White" />
	      <label for="white">White</label>
	      <input type="checkbox" id="blue" value="Blue" />
	      <label for="blue">Blue</label>
	      <input type="checkbox" id="black" value="Black" />
	      <label for="black">Black</label>
	      <input type="checkbox" id="red" value="Red" />
	      <label for="red">Red</label>
	      <input type="checkbox" id="green" value="Green" />
	      <label for="green">Green</label>
	      <input type="checkbox" id="colourless" value="Colourless" />
	      <label for="colourless">Colourless</label>
	    </fieldset>
	    <div>
	      <label for="colour-op">Operation</label>
	      <select id="colour-op" data-opcode>
	        <option value="exact">Exact match</option>
	        <option value="include">Include these colours</option>
	        <option value="at-most">At most these colours</option>
	      </select>
	    </div>
	    <div>
	      <input id="multicoloured" type="checkbox" data-additional-filter="colour:length>1" />
	      <label for="multicoloured">Multicoloured Only</label>
	    </div>
	  </div>`,
	  'rarity': `<div class="input-group" data-field-name="rarity.name">
	    <fieldset data-main-input>
	      <legend>Rarity</legend>
	      <input type="checkbox" id="core" value="Core" checked />
	      <label for="core">Core</label>
	      <input type="checkbox" id="signature" value="Signature" checked />
	      <label for="signature">Signature</label>
	      <input type="checkbox" id="token" value="Token" />
	      <label for="token">Token</label>
	      <input type="checkbox" id="common" value="Common" checked />
	      <label for="common">Common</label>
	      <input type="checkbox" id="rare" value="Rare" checked />
	      <label for="rare">Rare</label>
	      <input type="checkbox" id="epic" value="Epic" checked />
	      <label for="epic">Epic</label>
	      <input type="checkbox" id="mythic" value="Mythic" checked />
	      <label for="mythic">Mythic</label>
	    </fieldset>
	  </div>`,
	  'set': `<div class="input-group" data-field-name="set.name">
	    <fieldset data-main-input>
	      <legend>Set</legend>
	      <input type="checkbox" id="core_set" value="Core Set" />
	      <label for="core_set">Core Set</label>
	      <input type="checkbox" id="signatures" value="Signatures" />
	      <label for="signatures">Signatures</label>
	      <input type="checkbox" id="opening_ceremony" value="Opening Ceremony" />
	      <label for="opening_ceremony">Opening Ceremony</label>
	      <input type="checkbox" id="dd_icons" value="D&D: Icons" />
	      <label for="dd_icons">D&amp;D: Icons</label>
	      <input type="checkbox" id="helvault_unsealed" value="Helvault Unsealed" />
	      <label for="helvault_unsealed">Helvault Unsealed</label>
	    </fieldset>
	  </div>`,
	  'type': `<div class="input-group" data-field-name="type.name">
	    <fieldset data-main-input>
	      <legend>Type</legend>
	      <input type="checkbox" id="creature" value="Creature" />
	      <label for="creature">Creature</label>
	      <input type="checkbox" id="spell" value="Spell" />
	      <label for="spell">Spell</label>
	      <input type="checkbox" id="artifact" value="Artifact" />
	      <label for="artifact">Artifact</label>
	      <input type="checkbox" id="trap" value="Trap" />
	      <label for="trap">Trap</label>
	      <input type="checkbox" id="skill" value="Skill" />
	      <label for="skill">Skill</label>
	      <input type="checkbox" id="land" value="Land" />
	      <label for="land">Land</label>
	    </fieldset>
	  </div>`,
	  'subtype': `<div class="input-group" data-field-name="subtype.name">
	    <label for="subtype">Subtype</label>
	    <input type="text" id="subtype" data-main-input />
	  </div>`,
	  'chance': `<div class="input-group" data-field-name="chance">
	    <label for="chance">Chance</label>
	    <input type="number" id="chance" data-main-input />
	    <div>
	      <label for="chance-op">Operation</label>
	      <select id="chance-op" data-opcode>
	        <option value="=">=</option>
	        <option value="!=">!=</option>
	        <option value=">">&gt;</option>
	        <option value=">=">&gt;=</option>
	        <option value="<">&lt;</option>
	        <option value="<=">&lt;=</option>
	      </select>
	    </div>
		</div>`,
    'artist': `<div class="input-group" data-field-name="artist">
      <label for="artist">Artist</label>
      <input type="text" id="artist" data-main-input />
    </div>`,
    'legendary': `<div class="input-group" data-field-name="legendary">
      <label for="legendary">Legendary</label>
      <select id="legendary" data-main-input>
        <option value="">All</option>
        <option value="true">Legendary</option>
        <option value="false">Non-Legendary</option>
      </select>
    </div>`,
    'generates': `<div class="input-group" data-field-name="generates.name">
      <label for="generates">Generates</label>
      <input type="text" id="generates" data-main-input />
    </div>`,
    'reminders': `<div class="input-group" data-field-name="reminders.name">
      <label for="reminders">Keyword with Reminder Text</label>
      <input type="text" id="reminders" data-main-input />
    </div>`,
    'charges': `<div class="input-group" data-field-name="charges">
      <label for="charges">Charges</label>
      <input type="number" id="charges" data-main-input />
      <div>
        <label for="charges-op">Operation</label>
        <select id="charges-op" data-opcode>
          <option value="=">=</option>
          <option value="!=">!=</option>
          <option value=">">&gt;</option>
          <option value=">=">&gt;=</option>
          <option value="<">&lt;</option>
          <option value="<=">&lt;=</option>
        </select>
      </div>
    </div>`,
	}
}

customElements.define('ssrr-cardlist', SsrrCardList)

export default SsrrCardList
