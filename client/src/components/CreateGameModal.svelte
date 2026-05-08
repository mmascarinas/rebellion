<script lang="ts">
	let {
		onclose,
		oncreate
	}: {
		onclose: () => void
		oncreate: (config: {
			name: string
			maxPlayers: number
			isPrivate: boolean
			password: string
		}) => void
	} = $props()

	let name = $state('')
	let maxPlayers = $state(4)
	let isPrivate = $state(false)
	let password = $state('')

	const playerOptions = [2, 3, 4, 5, 6]

	const canCreate = $derived(
		name.trim().length > 0 && (!isPrivate || password.trim().length > 0)
	)

	function handleCreate() {
		if (!canCreate) return
		oncreate({
			name: name.trim(),
			maxPlayers,
			isPrivate,
			password: isPrivate ? password.trim() : ''
		})
	}

	function openDialog(node: HTMLDialogElement) {
		node.showModal()
	}

	function handleCancel(e: Event) {
		e.preventDefault()
		onclose()
	}

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) onclose()
	}
</script>

<dialog
	use:openDialog
	class="overlay"
	aria-label="Create Game"
	onclick={handleBackdropClick}
	oncancel={handleCancel}
>
	<div class="panel">
		<header class="panel-header">
			<h2 class="panel-title">Create Game</h2>
			<button class="close-btn" onclick={onclose}>✕</button>
		</header>

		<div class="form">
			<div class="field">
				<label class="label" for="game-name">Game Name</label>
				<input
					id="game-name"
					class="input"
					type="text"
					placeholder="Name your conspiracy..."
					maxlength="32"
					bind:value={name}
				/>
			</div>

			<div class="field">
				<span class="label">Max Players</span>
				<div class="player-select">
					{#each playerOptions as n}
						<button
							class="player-option"
							class:selected={maxPlayers === n}
							onclick={() => (maxPlayers = n)}
						>
							{n}
						</button>
					{/each}
				</div>
			</div>

			<div class="field">
				<span class="label">Visibility</span>
				<div class="toggle-row">
					<button
						class="toggle-option"
						class:selected={!isPrivate}
						onclick={() => (isPrivate = false)}
					>
						Public
					</button>
					<button
						class="toggle-option"
						class:selected={isPrivate}
						onclick={() => (isPrivate = true)}
					>
						Private
					</button>
				</div>
			</div>

			{#if isPrivate}
				<div class="field fade-in">
					<label class="label" for="game-password">Password</label>
					<input
						id="game-password"
						class="input"
						type="password"
						placeholder="Set a secret phrase..."
						maxlength="32"
						bind:value={password}
					/>
				</div>
			{/if}
		</div>

		<footer class="panel-footer">
			<button class="btn btn-cancel" onclick={onclose}>Cancel</button>
			<button
				class="btn btn-create"
				disabled={!canCreate}
				onclick={handleCreate}
			>
				Create
			</button>
		</footer>
	</div>
</dialog>

<style lang="scss">
	.overlay {
		position: fixed;
		inset: 0;
		border: none;
		background: transparent;
		width: 100%;
		height: 100%;
		max-width: 100%;
		max-height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
		animation: fadeIn 0.2s ease;
		padding: 24px;

		&::backdrop {
			background: rgba(0, 0, 0, 0.5);
			backdrop-filter: blur(4px);
		}
	}

	.panel {
		background: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: 12px;
		width: 100%;
		max-width: 420px;
		animation: slideUp 0.25s ease;
		overflow: hidden;
	}

	.panel-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 24px 28px 0;
	}

	.panel-title {
		font-size: 18px;
		font-weight: 600;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--color-text);
	}

	.close-btn {
		background: none;
		border: none;
		color: var(--color-text-secondary);
		font-size: 18px;
		cursor: pointer;
		padding: 4px;
		line-height: 1;
		transition: color 0.15s ease;

		&:hover {
			color: var(--color-text);
		}
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 24px;
		padding: 28px;
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.label {
		font-size: 11px;
		font-weight: 500;
		letter-spacing: 0.15em;
		text-transform: uppercase;
		color: var(--color-text-secondary);
	}

	.input {
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 6px;
		padding: 12px 14px;
		font-size: 14px;
		font-family: inherit;
		color: var(--color-text);
		outline: none;
		transition: border-color 0.2s ease;

		&::placeholder {
			color: var(--color-text-secondary);
			opacity: 0.5;
		}

		&:focus {
			border-color: var(--color-primary);
		}
	}

	.player-select {
		display: flex;
		gap: 8px;
	}

	.player-option {
		flex: 1;
		padding: 10px 0;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 6px;
		color: var(--color-text-secondary);
		font-size: 14px;
		font-weight: 500;
		font-family: inherit;
		cursor: pointer;
		transition: all 0.15s ease;

		&.selected {
			border-color: var(--color-primary);
			color: var(--color-primary);
			background: var(--color-bg);
		}
	}

	.toggle-row {
		display: flex;
		gap: 8px;
	}

	.toggle-option {
		flex: 1;
		padding: 10px 0;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 6px;
		color: var(--color-text-secondary);
		font-size: 13px;
		font-weight: 500;
		font-family: inherit;
		letter-spacing: 0.05em;
		text-transform: uppercase;
		cursor: pointer;
		transition: all 0.15s ease;

		&.selected {
			border-color: var(--color-primary);
			color: var(--color-primary);
			background: var(--color-bg);
		}
	}

	.panel-footer {
		display: flex;
		gap: 12px;
		padding: 0 28px 24px;
		justify-content: flex-end;
	}

	.btn {
		padding: 10px 24px;
		border-radius: 6px;
		font-size: 13px;
		font-weight: 500;
		font-family: inherit;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.btn-cancel {
		background: none;
		border: 1px solid var(--color-border);
		color: var(--color-text-secondary);

		&:hover {
			border-color: var(--color-text-secondary);
			color: var(--color-text);
		}
	}

	.btn-create {
		background: var(--color-primary);
		border: 1px solid var(--color-primary);
		color: #fff;

		&:hover:not(:disabled) {
			opacity: 0.9;
			transform: translateY(-1px);
		}

		&:disabled {
			opacity: 0.35;
			cursor: not-allowed;
		}
	}

	.fade-in {
		animation: fadeIn 0.2s ease;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(16px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@media (max-width: 330px) {
		.overlay {
			padding: 12px;
		}

		.panel-header {
			padding: 16px 16px 0;
		}

		.panel-title {
			font-size: 14px;
		}

		.form {
			padding: 16px;
			gap: 16px;
		}

		.input {
			padding: 10px 12px;
			font-size: 12px;
		}

		.player-option {
			padding: 8px 0;
			font-size: 12px;
		}

		.toggle-option {
			padding: 8px 0;
			font-size: 11px;
		}

		.panel-footer {
			padding: 0 16px 16px;
		}

		.btn {
			padding: 8px 16px;
			font-size: 11px;
		}
	}
</style>
