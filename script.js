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

    reset() {
        this.grid = Array.from({ length: ROWS }, () => Array(COLS).fill(null));
        this.activePuyoGroup = null;
        this.score = 0;
        this.nuisance = 0;
        this.pendingNuisance = 0;
        this.isAnimating = false;
        this.render();
        // Clear nuisance UI
        const nuisanceContainer = document.getElementById(`${this.playerId}-nuisance`);
        if (nuisanceContainer) nuisanceContainer.innerHTML = '';
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

        // Moves left for each player
        this.p1MovesLeft = 3;
        this.p2MovesLeft = 0;

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

        document.getElementById('p1-moves').innerText = this.p1MovesLeft;
        document.getElementById('p2-moves').innerText = this.p2MovesLeft;

        this.updateNuisanceUI(this.p1Board, 'p1-nuisance');
        this.updateNuisanceUI(this.p2Board, 'p2-nuisance');
    }

    updateNuisanceUI(board, elementId) {
        const container = document.getElementById(elementId);
        container.innerHTML = '';

        let amount = board.pendingNuisance;

        // Standard Puyo Nuisance Symbols
        const symbols = [
            { value: 720, className: 'nuisance-crown' },
            { value: 360, className: 'nuisance-moon' },
            { value: 180, className: 'nuisance-star' },
            { value: 30, className: 'nuisance-rock' },
            { value: 6, className: 'nuisance-big' },
            { value: 1, className: 'nuisance-small' }
        ];

        for (const sym of symbols) {
            while (amount >= sym.value) {
                const marker = document.createElement('div');
                marker.className = `nuisance-marker ${sym.className}`;
                container.appendChild(marker);
                amount -= sym.value;

                // Limit max markers to avoid overflow UI
                if (container.children.length > 10) return;
            }
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
            // 1. Apply gravity
            const dropped = await this.applyGravity(board);

            // 2. Find matches
            const matches = this.findMatches(board);

            if (matches.length > 0) {
                chainCount++;

                // 3. Calculate Score (Official Rule)
                const scoreData = this.calculateScore(matches, chainCount);
                board.score += scoreData.score;

                // 4. Calculate Nuisance
                // Standard rate: 70 points = 1 nuisance
                const nuisanceRate = 70;
                let nuisancePoints = scoreData.score + (board.nuisanceRemainder || 0);
                const nuisanceGenerated = Math.floor(nuisancePoints / nuisanceRate);
                board.nuisanceRemainder = nuisancePoints % nuisanceRate;

                board.nuisance += nuisanceGenerated;

                // 5. Remove matches
                await this.removeMatches(board, matches);
                somethingChanged = true;

                this.updateUI();
            } else {
                somethingChanged = dropped;
            }

            if (somethingChanged) {
                await new Promise(r => setTimeout(r, 300));
            }
        }

        let turnEnd = false;
        // Turn end logic
        if (this.turn === 'p1') {
            this.p1MovesLeft--;
            if (this.p1MovesLeft === 0) {
                turnEnd = true;
            }
        } else {
            this.p2MovesLeft--;
            if (this.p2MovesLeft === 0) {
                turnEnd = true;
            }
        }
        this.updateUI();

        if (chainCount > 0 || turnEnd) {
            // Calculate bonus moves for next player
            // Note: movesLeft has already been decremented for the current move
            let nextMovesBonus = chainCount * 2;
            if (this.turn === 'p1') {
                this.p2MovesLeft += nextMovesBonus;
            } else {
                this.p1MovesLeft += nextMovesBonus;
            }

            this.handleNuisance(board);

            if (this.checkGameOver()) {
                return;
            }

            this.switchTurn();
        } else {
            this.startTurn();
        }
    }

    checkGameOver() {
        // Check if either player has lost
        // Condition: Spawn point (0, 2) is blocked
        // We check both boards because garbage might have fallen? 
        // Actually, usually only the active player could have died from their own move or garbage falling on them.
        // But checking both is safe.

        if (this.p1Board.grid[0][2] || this.p1Board.grid[1][2]) {
            alert('Player 2 Wins!');
            this.reset();
            return true;
        }
        if (this.p2Board.grid[0][2] || this.p2Board.grid[1][2]) {
            alert('Player 1 Wins!');
            this.reset();
            return true;
        }
        return false;
    }

    reset() {
        this.p1Board.reset();
        this.p2Board.reset();
        this.turn = 'p1';
        this.p1MovesLeft = 3;
        this.p2MovesLeft = 0;
        this.p1Queue = [this.generateColors(), this.generateColors()];
        this.p2Queue = [this.generateColors(), this.generateColors()];
        this.updateUI();
        this.updateNextPuyoUI();
        this.startTurn();
    }

    calculateScore(matches, chainCount) {
        // Group matches by color and connectivity to determine bonuses
        // matches is a flat list of puyos. We need to reconstruct groups to calculate Group Bonus and Color Bonus.
        // Actually findMatches returns a flat list.
        // We should probably let findMatches return groups or re-group here.
        // For simplicity, let's re-group or modify findMatches.
        // But findMatches logic in this file returns a flat array of all matched puyos.
        // We can deduce colors easily.

        const uniqueColors = new Set(matches.map(p => p.color));
        const colorCount = uniqueColors.size;
        const puyoCount = matches.length;

        // Group Bonus is per group. We need to know the size of each connected group.
        // Since we don't have that info readily from flat list, let's approximate or re-calculate.
        // Re-calculating groups from the matches list:
        // We can just run a quick connected component search on the matches list?
        // Or better, modify findMatches to return groups.
        // Let's modify findMatches in a separate edit if needed, but for now let's assume we can get it.
        // Actually, let's just implement a helper here to group the matches.

        const groups = this.groupMatches(matches);

        // Constants (Puyo Puyo Tsu)
        const CHAIN_BONUS = [0, 8, 16, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 480, 512]; // 1-indexed (index 0 is chain 1)
        const COLOR_BONUS = [0, 0, 3, 6, 12, 24]; // 1, 2, 3, 4, 5 colors
        const GROUP_BONUS = [0, 0, 0, 0, 0, 2, 3, 4, 5, 6, 7, 10]; // 0-3, 4, 5, 6... 11+ is 10

        // Calculate Bonuses
        let cp = CHAIN_BONUS[Math.min(chainCount, CHAIN_BONUS.length) - 1] || 0;

        let cb = COLOR_BONUS[Math.min(colorCount, COLOR_BONUS.length) - 1] || 0;

        let gb = 0;
        for (const group of groups) {
            const size = group.length;
            if (size >= 11) gb += 10;
            else if (size >= 4) gb += GROUP_BONUS[size] || 0;
        }

        let totalBonus = cp + cb + gb;
        if (totalBonus === 0) totalBonus = 1; // Min 1
        if (totalBonus > 999) totalBonus = 999; // Max 999

        const score = (10 * puyoCount) * totalBonus;
        return { score };
    }

    groupMatches(matches) {
        // Helper to group puyos by connectivity
        // Since they are already matched, we just need to group them by color and adjacency?
        // Actually, findMatches logic already grouped them but flattened.
        // Let's just group by color for now as an approximation,
        // assuming matches of same color are one group (which is true 99% of time in simple chains).
        // BUT, separate groups of same color can exist (e.g. 4 red here, 4 red there).
        // For strict rules, we need adjacency.

        const groups = [];
        const visited = new Set();

        for (const p of matches) {
            if (visited.has(p)) continue;

            const group = [p];
            visited.add(p);
            const stack = [p];

            while (stack.length > 0) {
                const curr = stack.pop();
                // Check neighbors in matches list
                for (const other of matches) {
                    if (!visited.has(other) && other.color === curr.color) {
                        if (Math.abs(other.r - curr.r) + Math.abs(other.c - curr.c) === 1) {
                            visited.add(other);
                            group.push(other);
                            stack.push(other);
                        }
                    }
                }
            }
            groups.push(group);
        }
        return groups;
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
        // Usually nuisance falls at the START of your turn, or after you finish your chain, if you didn't offset everything.
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
        if (this.turn === 'p1') {
            this.p1MovesLeft += 3;
        } else {
            this.p2MovesLeft += 3;
        }
        this.updateUI();
        this.startTurn();
    }
}

window.onload = () => {
    const game = new Game();
};
