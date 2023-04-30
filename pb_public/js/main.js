// ---------------------------
// Toggle Nav
// ---------------------------

const swipeData = { x: null, y: null }

const expandMenu = localStorage.getItem('expand-menu')
const menuToggle = document.querySelector('#menu-toggle')
const siteNavigation = document.querySelector('#primary-navigation')
menuToggle.setAttribute('aria-expanded', expandMenu ?? true)
siteNavigation.setAttribute('data-expanded', expandMenu ?? true)

document.addEventListener("touchstart", (e) => {
  if (e.touches.length != 1) {
    swipeData.x = null
    swipeData.y = null
  }
  else if (e.touches[0].clientX < 50) {
    swipeData.x = e.touches[0].clientX
    swipeData.y = e.touches[0].clientY
  }
  else {
    const isOpened = menuToggle.getAttribute('aria-expanded') === 'true'
    
    if (isOpened) {
      swipeData.x = e.touches[0].clientX
      swipeData.y = e.touches[0].clientY
    }
    else {
      swipeData.x = null
      swipeData.y = null
    }
  }
}, false)

document.addEventListener("touchend", (e) => {
  swipeData.x = null
  swipeData.y = null
}, false)

document.addEventListener("touchcancel", (e) => {
  swipeData.x = null
  swipeData.y = null
}, false)

document.addEventListener("touchmove", (e) => {
  if (!swipeData.x || !swipeData.y) return

  e.preventDefault()

  const diffX = swipeData.x - e.touches[0].clientX
  const diffY = swipeData.y - e.touches[0].clientY

  // if X diff is more significant
  if (Math.abs(diffX) > Math.abs(diffY)) {
    toggleNav(diffX < 0) // toggle based on left/right
  }
}, false)

menuToggle.addEventListener('click', () => {
  const isOpened = menuToggle.getAttribute('aria-expanded') === 'true'
  toggleNav(!isOpened)
})

function toggleNav(open) {
  menuToggle.setAttribute('aria-expanded', open)
  siteNavigation.setAttribute('data-expanded', open)
  localStorage.setItem('expand-menu', open)
}

// ---------------------------
// Modals
// ---------------------------

function confirmModal({ titleText, bodyText, onConfirm, onCancel }) {
  openModal('#confirm-modal', titleText, bodyText, onConfirm, onCancel)
}

function alertModal({ titleText, bodyText, onConfirm }) {
  openModal('#alert-modal', titleText, bodyText, onConfirm)
}

function openModal(modalSelector, titleText, bodyText, onConfirm, onCancel) {
  const modal = document.querySelector(modalSelector)
  const title = modal.querySelector('h2')
  const body = modal.querySelector('p')
  const confirmButton = modal.querySelector('.confirm-button')
  const cancelButton = modal.querySelector('.cancel-button')

  title.textContent = titleText
  body.innerHTML = bodyText
  confirmButton.onclick = () => {
    onConfirm()
    modal.close()
  }
  if (cancelButton) {
    cancelButton.onclick = () => {
      onCancel()
      modal.close()
    }
  }

  modal.showModal()
}

// ---------------------------
// Errors
// ---------------------------

function alertOnError(func) {
  try {
    func()
  }
  catch (e) {
    alertModal({
      titleText: 'An error occurred',
      bodyText: formatErrorMessage(e),
      onConfirm: () => { },
    })
  }
}

async function alertOnErrorAsync(func) {
  try {
    await func()
  }
  catch (e) {
    alertModal({
      titleText: 'An error occurred',
      bodyText: formatErrorMessage(e),
      onConfirm: () => { },
    })
  }
}

function formatErrorMessage(e) {
  if (!e.data?.message) return JSON.stringify(e.message || e.data || e)

  let message = e.data.message
  if (e.data.data) {
    message += '<br>' + Object.entries(e.data.data).map(([key, data]) => `${key}: ${data.message}`).join('<br>')
  }
  return message
}

// ---------------------------
// Outage Notification
// ---------------------------

const latestNotificationSeen = parseInt(localStorage.getItem('latest-notification-seen') ?? '0')
const currentNotification = 1
const outageTime = new Date('2023-01-01T00:00:00.000-06:00')

if (latestNotificationSeen < currentNotification && outageTime > new Date()) {
  alertModal({
    titleText: 'Upcoming Outage',
    bodyText: `The site will be offline to release new features for approximately 30 minutes starting at ${outageTime}`,
    onConfirm: () => {
      localStorage.setItem('latest-notification-seen', currentNotification)
    },
  })
}

export { alertModal, confirmModal, alertOnError, alertOnErrorAsync }