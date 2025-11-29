const ROWS = 12;
const COLS = 6;
const COLORS = ['red', 'green', 'blue', 'yellow'];

class Puyo {
    constructor(color, r, c) {
        this.color = color;
        this.r = r;
        this.c = c;
        this.element = null;
    }
}

class Board {
    constructor(playerId) {
        this.playerId = playerId;
        this.grid = Array.from({ length: ROWS }, () => Array(COLS).fill(null));
        this.element = document.getElementById(`${playerId}-grid`);
        this.activePuyoGroup = null; // { puyos: [p1, p2], r, c, rotation }
        this.score = 0;
        this.nuisance = 0; // Generated this turn
        this.pendingNuisance = 0; // Incoming from opponent
        this.isAnimating = false;
    }

    render() {
        this.element.innerHTML = '';
        // Render static grid
        for (let r = 0; r < ROWS; r++) {
            for (let c = 0; c < COLS; c++) {
                const cell = document.createElement('div');
                cell.className = 'cell';
                this.element.appendChild(cell);

                const puyo = this.grid[r][c];
                if (puyo) {
                    this.renderPuyo(puyo);
                }
            }
        }
        // Render active puyo
        if (this.activePuyoGroup) {
            this.activePuyoGroup.puyos.forEach(p => this.renderPuyo(p));
        }
    }

    renderPuyo(puyo) {
        const el = document.createElement('div');
        el.className = `puyo ${puyo.color}`;
        el.style.top = `${puyo.r * 30}px`;
        el.style.left = `${puyo.c * 30}px`;
        this.element.appendChild(el);
        puyo.element = el;
    }

    spawnPuyo(colors) {
        // Spawn at top center (0, 2) and (1, 2) usually, but let's say (0, 2) is pivot
        // Simplified: Pivot at (0, 2), other at (-1, 2) initially then drops?
        // Let's spawn at (0, 2) and (1, 2) for simplicity if empty
        if (this.grid[0][2] || this.grid[1][2]) {
            return false; // Game Over
        }

        const p1 = new Puyo(colors[0], 1, 2); // Main
        const p2 = new Puyo(colors[1], 0, 2); // Sub (above)

        this.activePuyoGroup = {
            puyos: [p1, p2],
            rotation: 0 // 0: up, 1: right, 2: down, 3: left (relative to main)
        };
        this.render();
        return true;
    }
}

class Game {
    constructor() {
        this.p1Board = new Board('p1');
        this.p2Board = new Board('p2');
        this.turn = 'p1'; // 'p1' or 'p2'
        this.movesLeft = 3;

        // Independent queues
        this.p1Queue = [this.generateColors(), this.generateColors()];
        this.p2Queue = [this.generateColors(), this.generateColors()];

        this.updateUI();
        this.updateNextPuyoUI();
        this.startTurn();

        document.addEventListener('keydown', (e) => this.handleInput(e));
    }

    generateColors() {
        return [
            COLORS[Math.floor(Math.random() * COLORS.length)],
            COLORS[Math.floor(Math.random() * COLORS.length)]
        ];
    }

    startTurn() {
        const board = this.turn === 'p1' ? this.p1Board : this.p2Board;
        const queue = this.turn === 'p1' ? this.p1Queue : this.p2Queue;

        const currentColors = queue.shift();
        queue.push(this.generateColors());

        if (!board.spawnPuyo(currentColors)) {
            alert(`${this.turn === 'p1' ? 'Player 2' : 'Player 1'} Wins!`);
            return;
        }
        this.updateNextPuyoUI();
    }

    updateNextPuyoUI() {
        this.renderQueue(this.p1Queue, 'p1');
        this.renderQueue(this.p2Queue, 'p2');
    }

    renderQueue(queue, playerId) {
        const nextEl = document.getElementById(`${playerId}-next`);
        const nextNextEl = document.getElementById(`${playerId}-next-next`);

        nextEl.innerHTML = '';
        nextNextEl.innerHTML = '';

        const [nextColors, nextNextColors] = queue;

        // Render Next
        // colors[0] is Main (Bottom), colors[1] is Sub (Top)
        // We want Top puyo in Row 1, Bottom puyo in Row 2
        const p1_sub = document.createElement('div');
        p1_sub.className = `puyo ${nextColors[1]}`;
        nextEl.appendChild(p1_sub);
        const p1_main = document.createElement('div');
        p1_main.className = `puyo ${nextColors[0]}`;
        nextEl.appendChild(p1_main);

        // Render Next-Next
        const nn1_sub = document.createElement('div');
        nn1_sub.className = `puyo ${nextNextColors[1]}`;
        nextNextEl.appendChild(nn1_sub);
        const nn1_main = document.createElement('div');
        nn1_main.className = `puyo ${nextNextColors[0]}`;
        nextNextEl.appendChild(nn1_main);
    }

    updateUI() {
        document.getElementById('turn-indicator').innerText = `Turn: ${this.turn === 'p1' ? 'Player 1' : 'Player 2'}`;
        document.getElementById('p1-moves').innerText = this.turn === 'p1' ? 3 - this.movesLeft : 0; // Wait, logic is "moves done" or "moves left"? User said "3 moves".
        // Let's show moves remaining or moves used. User said "3 moves operation".
        // Let's show moves remaining.
        document.getElementById('p1-moves').innerText = this.turn === 'p1' ? this.movesLeft : '-';
        document.getElementById('p2-moves').innerText = this.turn === 'p2' ? this.movesLeft : '-';

        this.updateNuisanceUI(this.p1Board, 'p1-nuisance');
        this.updateNuisanceUI(this.p2Board, 'p2-nuisance');
    }

    updateNuisanceUI(board, elementId) {
        const container = document.getElementById(elementId);
        container.innerHTML = '';
        // 1 marker per 1 nuisance (simplified)
        // Max 30 markers to avoid overflow
        const count = Math.min(board.pendingNuisance, 30);
        for (let i = 0; i < count; i++) {
            const marker = document.createElement('div');
            marker.className = 'nuisance-marker';
            container.appendChild(marker);
        }
    }

    handleInput(e) {
        const board = this.turn === 'p1' ? this.p1Board : this.p2Board;
        if (!board.activePuyoGroup || board.isAnimating) return;

        // Shared controls for PoC
        switch (e.key) {
            case 'ArrowLeft':
                this.move(-1);
                break;
            case 'ArrowRight':
                this.move(1);
                break;
            case 'ArrowDown':
                this.drop();
                break;
            case 'z':
                this.rotate(-1);
                break;
            case 'x':
                this.rotate(1);
                break;
        }
    }

    move(dir) {
        const board = this.turn === 'p1' ? this.p1Board : this.p2Board;
        const group = board.activePuyoGroup;

        // Check collision
        let canMove = true;
        for (const p of group.puyos) {
            const newC = p.c + dir;
            if (newC < 0 || newC >= COLS || board.grid[p.r][newC]) {
                canMove = false;
                break;
            }
        }

        if (canMove) {
            group.puyos.forEach(p => p.c += dir);
            board.render();
        }
    }

    rotate(dir) {
        // Simple rotation logic placeholder
        // Need to implement rotation system (around pivot)
        const board = this.turn === 'p1' ? this.p1Board : this.p2Board;
        const group = board.activePuyoGroup;
        const main = group.puyos[0];
        const sub = group.puyos[1];

        // Current relative pos
        const dr = sub.r - main.r;
        const dc = sub.c - main.c;

        // Rotate (90 deg)
        // (dr, dc) -> (dc, -dr) for CW (x)
        // (dr, dc) -> (-dc, dr) for CCW (z)

        let newDr, newDc;
        if (dir === 1) { // CW
            newDr = dc;
            newDc = -dr;
        } else { // CCW
            newDr = -dc;
            newDc = dr;
        }

        const newR = main.r + newDr;
        const newC = main.c + newDc;

        // Check bounds/collision
        if (newC >= 0 && newC < COLS && newR >= 0 && newR < ROWS && !board.grid[newR][newC]) {
            sub.r = newR;
            sub.c = newC;
            board.render();
        } else {
            // Wall kick? For PoC maybe just block
        }
    }

    drop() {
        const board = this.turn === 'p1' ? this.p1Board : this.p2Board;
        const group = board.activePuyoGroup;

        // Hard drop or soft drop? "3 moves or fire". 
        // Usually "Down" moves one step. "Up" or separate button for hard drop.
        // Let's make Down move one step, if blocked, lock.

        let canDrop = true;
        for (const p of group.puyos) {
            const newR = p.r + 1;
            if (newR >= ROWS || board.grid[newR][p.c]) {
                canDrop = false;
                break;
            }
        }

        if (canDrop) {
            group.puyos.forEach(p => p.r += 1);
            board.render();
        } else {
            this.lockPuyos();
        }
    }

    lockPuyos() {
        const board = this.turn === 'p1' ? this.p1Board : this.p2Board;
        const group = board.activePuyoGroup;

        // Place in grid
        group.puyos.forEach(p => {
            board.grid[p.r][p.c] = p;
        });
        board.activePuyoGroup = null;
        board.render();

        // Check for matches / gravity
        // This will be async/animated in real game
        this.resolveMatches(board);
    }

    async resolveMatches(board) {
        let chainCount = 0;
        let somethingChanged = true;

        while (somethingChanged) {
            // 1. Apply gravity (drop floating puyos)
            const dropped = await this.applyGravity(board);

            // 2. Find matches
            const matches = this.findMatches(board);

            if (matches.length > 0) {
                // 3. Remove matches
                await this.removeMatches(board, matches);
                chainCount++;
                somethingChanged = true;

                // Calculate score/garbage (simplified)
                const points = matches.length * 10 * chainCount; // Basic scoring
                board.score += points;

                // Add nuisance to opponent (simplified: 1 nuisance per 4 puyos cleared)
                // In real Puyo, it's based on score.
                // Let's stick to the user's request: "Player 1 moves 3 times OR fires -> nuisance calc"
                // Actually, nuisance is calculated per chain but sent at end of turn?
                // Or sent immediately? User said "Player 1 moves 3 times or fires... nuisance calc... switch".
                // So we accumulate nuisance during the chain.

                // For PoC, let's just track score for now, and calculate nuisance at end of turn.
                // Or maybe calculate nuisance here and store in a buffer.
                const nuisanceGenerated = Math.floor(points / 100); // Dummy formula
                board.nuisance += nuisanceGenerated;

                this.updateUI();
            } else {
                somethingChanged = dropped; // If only dropped but no matches, we might need to check matches again? 
                // Actually, if dropped, we loop again to check matches.
                // If nothing dropped and no matches, we are stable.
            }

            if (somethingChanged) {
                await new Promise(r => setTimeout(r, 300)); // Animation delay
            }
        }

        // Turn end logic
        this.movesLeft--;
        this.updateUI();

        // Switch turn if chain occurred OR moves ran out
        if (chainCount > 0 || this.movesLeft <= 0) {
            // Handle nuisance transfer before switching
            this.handleNuisance(board);
            this.switchTurn();
        } else {
            this.startTurn();
        }
    }

    async applyGravity(board) {
        let dropped = false;
        // Process column by column
        for (let c = 0; c < COLS; c++) {
            let writeIdx = ROWS - 1;
            for (let r = ROWS - 1; r >= 0; r--) {
                if (board.grid[r][c]) {
                    if (r !== writeIdx) {
                        board.grid[writeIdx][c] = board.grid[r][c];
                        board.grid[r][c] = null;
                        board.grid[writeIdx][c].r = writeIdx;
                        dropped = true;
                    }
                    writeIdx--;
                }
            }
        }
        board.render();
        if (dropped) await new Promise(r => setTimeout(r, 200));
        return dropped;
    }

    findMatches(board) {
        const matches = [];
        const visited = Array.from({ length: ROWS }, () => Array(COLS).fill(false));

        for (let r = 0; r < ROWS; r++) {
            for (let c = 0; c < COLS; c++) {
                if (board.grid[r][c] && !visited[r][c] && board.grid[r][c].color !== 'garbage') {
                    const group = this.getConnectedGroup(board, r, c, board.grid[r][c].color, visited);
                    if (group.length >= 4) {
                        matches.push(...group);
                    }
                }
            }
        }
        return matches;
    }

    getConnectedGroup(board, r, c, color, visited) {
        const group = [];
        const stack = [{ r, c }];
        visited[r][c] = true;
        group.push(board.grid[r][c]);

        while (stack.length > 0) {
            const curr = stack.pop();
            const dirs = [[0, 1], [0, -1], [1, 0], [-1, 0]];

            for (const [dr, dc] of dirs) {
                const nr = curr.r + dr;
                const nc = curr.c + dc;

                if (nr >= 0 && nr < ROWS && nc >= 0 && nc < COLS &&
                    board.grid[nr][nc] && !visited[nr][nc] &&
                    board.grid[nr][nc].color === color) {
                    visited[nr][nc] = true;
                    group.push(board.grid[nr][nc]);
                    stack.push({ r: nr, c: nc });
                }
            }
        }
        return group;
    }

    async removeMatches(board, matches) {
        matches.forEach(p => {
            board.grid[p.r][p.c] = null;
            // Also remove adjacent garbage
            const dirs = [[0, 1], [0, -1], [1, 0], [-1, 0]];
            dirs.forEach(([dr, dc]) => {
                const nr = p.r + dr;
                const nc = p.c + dc;
                if (nr >= 0 && nr < ROWS && nc >= 0 && nc < COLS &&
                    board.grid[nr][nc] && board.grid[nr][nc].color === 'garbage') {
                    board.grid[nr][nc] = null;
                }
            });
        });
        board.render();
        await new Promise(r => setTimeout(r, 300));
    }

    handleNuisance(activeBoard) {
        // Calculate nuisance sent to opponent
        // In this PoC, activeBoard.nuisance is what they generated this turn
        const opponentBoard = activeBoard.playerId === 'p1' ? this.p2Board : this.p1Board;

        // Offset logic
        // If opponent has pending nuisance, offset it
        // Note: In real Puyo, nuisance is in a "tray".
        // Let's assume board.nuisance is the "tray" of incoming nuisance?
        // Wait, in my Board class, `nuisance` was initialized to 0.
        // Let's redefine:
        // `board.pendingNuisance`: Nuisance waiting to fall on this board.
        // `board.generatedNuisance`: Nuisance generated by this board in current turn.

        // Let's use `board.nuisance` as "Pending Nuisance to fall on me".
        // And we calculate `generated` locally.

        const generated = activeBoard.nuisance; // I used this field to store generated points in resolveMatches
        activeBoard.nuisance = 0; // Reset generated count for next turn

        if (generated > 0) {
            if (activeBoard.pendingNuisance > 0) {
                // Offset own pending nuisance
                if (generated >= activeBoard.pendingNuisance) {
                    const remainder = generated - activeBoard.pendingNuisance;
                    activeBoard.pendingNuisance = 0;
                    opponentBoard.pendingNuisance = (opponentBoard.pendingNuisance || 0) + remainder;
                } else {
                    activeBoard.pendingNuisance -= generated;
                }
            } else {
                // Send all to opponent
                opponentBoard.pendingNuisance = (opponentBoard.pendingNuisance || 0) + generated;
            }
        }

        // Drop nuisance if any remains on active player?
        // Rules: Nuisance falls on the player whose turn just ENDED? 
        // No, nuisance falls on the player who is ABOUT TO START or AFTER they finish?
        // User said: "(1) p1 moves... nuisance calc -> p2 turn. (2) p2 moves... nuisance calc -> offset -> p1 turn".
        // Usually nuisance falls at the START of your turn, or after you fail to offset.
        // Standard Puyo: Garbage falls after you finish your chain, if you didn't offset everything.
        // So if P1 finishes, and has pending garbage, it falls NOW.

        if (activeBoard.pendingNuisance > 0) {
            this.dropNuisance(activeBoard);
        }

        this.updateUI();
    }

    dropNuisance(board) {
        // Drop garbage puyos from top
        // Simplified: Random columns
        let count = board.pendingNuisance;
        board.pendingNuisance = 0; // Clear after dropping

        // Limit max garbage per turn to avoid infinite loops or crashes
        count = Math.min(count, 30);

        for (let i = 0; i < count; i++) {
            const c = Math.floor(Math.random() * COLS);
            // Find drop position
            // We just put it at top and let gravity handle it next frame?
            // Or force it into grid?
            // Let's force it into the first available spot from top
            let r = 0;
            while (r < ROWS && !board.grid[r][c]) {
                r++;
            }
            r--; // One above

            if (r >= 0) {
                board.grid[r][c] = new Puyo('garbage', r, c);
            }
        }
        board.render();
    }

    switchTurn() {
        this.turn = this.turn === 'p1' ? 'p2' : 'p1';
        this.movesLeft = 3;
        this.updateUI();
        this.startTurn();
    }
}

window.onload = () => {
    const game = new Game();
};
