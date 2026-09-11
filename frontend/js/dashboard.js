document.addEventListener('DOMContentLoaded', async () => {
    // check if user is logged in
    try {
        const user = await request('/me')
        document.getElementById('welcome').textContent = `Welcome, ${user.email}`
    } catch (err) {
        // not logged in, redirect back to login
        window.location.href = '/index.html'
        return
    }

    document.getElementById('logout-btn').addEventListener('click', logout)
    document.getElementById('start-btn').addEventListener('click', startChallenge)
})

async function logout() {
    try {
        await request('/logout', { method: 'POST' })
    } finally {
        window.location.href = '/index.html'
    }
}

function startChallenge() {
    const category = document.getElementById('category').value
    const level = document.getElementById('level').value
    const amount = document.getElementById('amount').value

        sessionStorage.setItem('challenge', JSON.stringify({ category, level, amount }))
    window.location.href = '/challenge.html'
}