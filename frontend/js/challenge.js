let questions = []
let currentIndex = 0
let correct = 0

document.addEventListener('DOMContentLoaded', async () => {
const stored = sessionStorage.getItem('challenge')
    if (!stored) {
        window.location.href = '/dashboard.html'
        return
    }

    const { category, level, amount } = JSON.parse(stored)
    console.log('params:', category, level, amount)

    try {
        const data = await request('/challenges/generate', {
            method: 'POST',
            body: JSON.stringify({ category, level, amount: parseInt(amount) })
        })
        questions = data.questions
        showQuestion()
    } catch (err) {
        document.getElementById('character').textContent = err.message
    }

    document.getElementById('restart-btn').addEventListener('click', () => {
        window.location.href = '/dashboard.html'
    })
})

function showQuestion() {
    if (currentIndex >= questions.length) {
        showResults()
        return
    }

    const q = questions[currentIndex]
    document.getElementById('character').textContent = q.character
    document.getElementById('feedback').textContent = ''

    const optionsEl = document.getElementById('options')
    optionsEl.innerHTML = ''

    q.options.forEach(option => {
        const btn = document.createElement('button')
        btn.textContent = option
        btn.addEventListener('click', () => handleAnswer(option, q.answer))
        optionsEl.appendChild(btn)
    })
}

function handleAnswer(selected, answer) {
    const feedbackEl = document.getElementById('feedback')

    if (selected === answer) {
        correct++
        feedbackEl.style.color = 'green'
        feedbackEl.textContent = 'Correct!'
    } else {
        feedbackEl.style.color = 'red'
        feedbackEl.textContent = `Wrong. Answer was: ${answer}`
    }

    // move to next question after a short delay
    setTimeout(() => {
        currentIndex++
        showQuestion()
    }, 1000)
}

function showResults() {
    document.getElementById('question-section').style.display = 'none'
    document.getElementById('results-section').style.display = 'block'
    document.getElementById('score').textContent = 
        `You got ${correct} out of ${questions.length} correct.`
}