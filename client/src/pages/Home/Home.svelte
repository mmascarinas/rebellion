<script lang="ts">
	import { push } from '../../lib/routing.svelte'
	import CreateGameButton from './CreateGameButton.svelte'

	let hovered = $state<string | null>(null)
	let showCreateGame = $state(false)
</script>

<div class="home">
	<div class="hero">
		<h1 class="title">REBELLION</h1>
		<p class="tagline">Deceive. Manipulate. Survive.</p>
	</div>

	<div class="actions">
		<button
			class="action-btn"
			class:active={hovered === 'create'}
			onmouseenter={() => (hovered = 'create')}
			onmouseleave={() => (hovered = null)}
			onclick={() => (showCreateGame = true)}
		>
			<span class="btn-icon">+</span>
			<span class="btn-label">Create Game</span>
		</button>

		<button
			class="action-btn"
			class:active={hovered === 'find'}
			onmouseenter={() => (hovered = 'find')}
			onmouseleave={() => (hovered = null)}
			onclick={() => push('/lobby')}
		>
			<span class="btn-icon">⌕</span>
			<span class="btn-label">Find Games</span>
		</button>

		<button
			class="action-btn"
			class:active={hovered === 'join'}
			onmouseenter={() => (hovered = 'join')}
			onmouseleave={() => (hovered = null)}
		>
			<span class="btn-icon">→</span>
			<span class="btn-label">Join Game</span>
		</button>
	</div>

	<p class="footer-note">The last one standing wins.</p>
</div>

{#if showCreateGame}
	<CreateGameButton
		onclose={() => (showCreateGame = false)}
		oncreate={(config) => {
			console.log('Create game:', config)
			showCreateGame = false
		}}
	/>
{/if}

<style lang="scss">
	.home {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		gap: 64px;
		padding: 48px 24px;
	}

	.hero {
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
	}

	.title {
		font-size: clamp(32px, 10vw, 80px);
		font-weight: 700;
		letter-spacing: 0.25em;
		color: var(--color-text);
		position: relative;

		&::after {
			content: '';
			display: block;
			width: 48px;
			height: 2px;
			background: var(--color-primary);
			margin: 20px auto 0;
		}
	}

	.tagline {
		font-size: 14px;
		font-weight: 300;
		letter-spacing: 0.35em;
		text-transform: uppercase;
		color: var(--color-text-secondary);
		margin-top: 8px;
	}

	.actions {
		display: flex;
		gap: 24px;
		flex-wrap: wrap;
		justify-content: center;
	}

	.action-btn {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 14px;
		width: 160px;
		padding: 32px 24px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 8px;
		color: var(--color-text);
		cursor: pointer;
		transition: all 0.25s ease;

		&:hover,
		&.active {
			border-color: var(--color-primary);
			transform: translateY(-4px);
			box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
		}

		.btn-icon {
			font-size: 28px;
			font-weight: 300;
			color: var(--color-primary);
			line-height: 1;
		}

		.btn-label {
			font-size: 13px;
			font-weight: 500;
			letter-spacing: 0.1em;
			text-transform: uppercase;
		}
	}

	.footer-note {
		font-size: 12px;
		font-weight: 300;
		letter-spacing: 0.2em;
		text-align: center;
		text-transform: uppercase;
		color: var(--color-text-secondary);
		opacity: 0.7;
	}
</style>
