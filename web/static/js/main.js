function loadResizer() {
    const canvas = document.querySelector("canvas#content")
    if (!canvas) {
        return
    }
    const observer = new ResizeObserver(() => {
        canvas.setAttribute("width", canvas.clientWidth)
        canvas.setAttribute("height", canvas.clientHeight)
    })
    observer.observe(canvas)
}

function updateState(cells) {
    const cellSize = 20
    const canvas = document.querySelector("canvas#content")
    const ctx = canvas.getContext("2d")
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    ctx.save()
    cells.forEach((cell) => {
        ctx.fillRect((canvas.width / 2) + (cell[1] * cellSize), (canvas.height / 2) + (cell[0] * cellSize), cellSize, cellSize)
    })
    ctx.restore()
}

function setupWS() {
    let socket

    setInterval(() => {
        socket?.send("")
    }, 1000)

    document.addEventListener('htmx:wsOpen', (e) => {
        socket = e.detail.socketWrapper;
    });
}

setupWS()
loadResizer && loadResizer()
window.onload = loadResizer
document.addEventListener("htmx:afterSwap", () => loadResizer())

document.addEventListener("htmx:wsAfterMessage", (event) => {
    const { action, payload } = JSON.parse(event.detail.message)
    switch (action) {
        case "updateState":
            updateState(payload)
            return
        case "pauseToggled":
            {
                const togglePauseButton = document.getElementById("togglePause")
                togglePauseButton.innerText = payload ? "Pause" : "Play"
            }
    }
})
