document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('register-btn').addEventListener('click', register)
    document.getElementById('login-btn').addEventListener('click', login)
})

async function register() {
    const email = document.getElementById('register-email').value
    const password = document.getElementById('register-password').value
    const errorEl = document.getElementById('register-error')

    try {
        await request('/register', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        })
        // show success message instead of redirecting
        errorEl.style.color = 'green'
        errorEl.textContent = 'Account created. Please login.'
    } catch (err) {
        errorEl.style.color = 'red'
        errorEl.textContent = err.message
    }
}

async function login() {
    const email = document.getElementById('login-email').value
    const password = document.getElementById('login-password').value
    const errorEl = document.getElementById('login-error')

    try {
        await request('/login', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        })
        window.location.href = '/dashboard.html'
    } catch (err) {
        errorEl.textContent = err.message
    }
}