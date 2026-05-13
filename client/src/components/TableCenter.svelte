<script lang="ts">
	interface Props {
		deckRemaining: number
		treasury: number
	}

	let { deckRemaining, treasury }: Props = $props()
</script>

<div class="table-center">
	<div class="deck-stack">
		<div class="deck-card"></div>
		<div class="deck-card"></div>
		<div class="deck-card top"></div>
		<span class="deck-count">{deckRemaining}</span>
	</div>
	<div class="treasury">
		<span class="treasury-icon">🪙</span>
		<span class="treasury-amount">{treasury}</span>
	</div>
</div>

<style lang="scss">
	.table-center {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		z-index: 2;
	}

	.deck-stack {
		position: relative;
		width: 72px;
		height: 100px;
		filter: drop-shadow(
			0 8px 24px
				color-mix(in srgb, var(--color-primary) 20%, rgba(0, 0, 0, 0.3))
		);
	}

	.deck-card {
		position: absolute;
		inset: 0;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: 8px;

		&:nth-child(1) {
			transform: translate(3px, 3px) rotate(2deg);
			opacity: 0.5;
		}
		&:nth-child(2) {
			transform: translate(1.5px, 1.5px) rotate(-1deg);
			opacity: 0.7;
		}
		&.top {
			transform: none;
			background: linear-gradient(
				135deg,
				var(--color-surface),
				color-mix(in srgb, var(--color-primary) 25%, var(--color-surface))
			);
			border-color: color-mix(
				in srgb,
				var(--color-primary) 30%,
				var(--color-border)
			);
			box-shadow:
				0 4px 16px
					color-mix(in srgb, var(--color-primary) 20%, rgba(0, 0, 0, 0.2)),
				inset 0 1px 0
					color-mix(in srgb, var(--color-primary) 15%, rgba(255, 255, 255, 0.1));
		}
	}

	.deck-count {
		position: absolute;
		bottom: -24px;
		left: 50%;
		transform: translateX(-50%);
		font-size: 11px;
		font-weight: 500;
		color: var(--color-text-secondary);
	}

	.treasury {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 8px 18px;
		background: linear-gradient(
			135deg,
			var(--color-surface),
			color-mix(in srgb, #ffb800 6%, var(--color-surface))
		);
		border: 1px solid color-mix(in srgb, #ffb800 20%, var(--color-border));
		border-radius: 20px;
		margin-top: 12px;
		box-shadow:
			0 4px 16px color-mix(in srgb, #ffb800 10%, transparent),
			inset 0 1px 0 color-mix(in srgb, #ffb800 8%, rgba(255, 255, 255, 0.1));
	}

	.treasury-icon {
		font-size: 16px;
		line-height: 1;
	}

	.treasury-amount {
		font-size: 14px;
		font-weight: 700;
		color: var(--color-text);
		letter-spacing: 0.02em;
	}
</style>
