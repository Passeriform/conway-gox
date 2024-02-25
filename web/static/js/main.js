function loadResizer() {
    const canvas = document.getElementById("game-map")
    const observer = new ResizeObserver(() => {
        canvas.setAttribute("width", canvas.clientWidth)
        canvas.setAttribute("height", canvas.clientHeight)
    })
    observer.observe(canvas)
}

function updateState(cells) {
    const cellSize = 20
    const canvas = document.getElementById("game-map")
    const ctx = canvas.getContext("2d")
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    ctx.save()
    cells.forEach((cell) => {
        ctx.fillRect(cell[1] * cellSize, cell[0] * cellSize, cellSize, cellSize)
    })
    ctx.restore()
}

window.onload = loadResizer

document.addEventListener("htmx:wsAfterMessage", (event) => {
    updateState(JSON.parse(event.detail.message))
})