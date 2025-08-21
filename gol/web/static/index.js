const gridCanvas = document.getElementById("grid");
const ctx = gridCanvas.getContext("2d");

const ws = new WebSocket(wsPath);
ws.onmessage = (event) => {
	const grid = JSON.parse(event.data);
	for (var y = 0; y < grid.length; y++) {
		for (var x = 0; x < grid[y].length; x++) {
			ctx.fillStyle = grid[y][x] == 1 ? "black" : "white";
			ctx.fillRect(10 * x, 10 * y, 10, 10)
		}
	}
}
