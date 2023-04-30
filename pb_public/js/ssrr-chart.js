class SsrrChart extends HTMLElement {
  constructor() {
    super()
    const shadow = this.attachShadow({ mode: 'open' })
    this.canvas = document.createElement('canvas')
    shadow.appendChild(this.canvas)
  }

  render({ type, title, data }) {
    if (!this.chart) {
      this.chart = new Chart(
        this.canvas,
        { 
          type: type,
          data: {},
          options: {
            plugins: {
              legend: {
                display: false,
              },
              title: {
                display: true,
                text: title,
              },
            },
          },
        },
      )
    }

    this.chart.data = {
      labels: data.map((d) => d.label),
      datasets: [{
        backgroundColor: data.map((d) => d.colour),
        data: data.map((d) => d.value),
      }],
    }

    this.chart.update()
  }
}

customElements.define('ssrr-chart', SsrrChart);

export default SsrrChart
