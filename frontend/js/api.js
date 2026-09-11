const API = 'http://localhost:8080/api'

async function request(path, options = {}) {
    const res = await fetch(`${API}${path}`, {
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        ...options
    })

    const data = await res.json()

    if (!res.ok) throw new Error(data.error || 'something went wrong')

    return data
}