<script lang="ts">
	type CardRole = 'courtesan' | 'mercenary' | 'baron' | 'consul' | 'admiral'

	const roleLabels: Record<CardRole, string> = {
		courtesan: 'Courtesan',
		mercenary: 'Mercenary',
		baron: 'Baron',
		consul: 'Consul',
		admiral: 'Admiral'
	}

	const roleIcons: Record<CardRole, string> = {
		courtesan: '♛',
		mercenary: '⚔',
		baron: '♜',
		consul: '⚖',
		admiral: '⚓'
	}

	interface Props {
		role: CardRole
		revealed?: boolean
		isOwn?: boolean
		facedown?: boolean
	}

	let {
		role,
		revealed = false,
		isOwn = false,
		facedown = false
	}: Props = $props()
</script>

{#if facedown}
	<div class="card facedown">
		<span class="card-back-pattern">⟡</span>
	</div>
{:else}
	<div class="card role-{role}" class:revealed class:own={isOwn}>
		<span class="card-icon">{roleIcons[role]}</span>
		<span class="card-label">{roleLabels[role]}</span>
	</div>
{/if}

<style lang="scss">
	.card {
		width: 48px;
		height: 68px;
		border-radius: 6px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 2px;
		font-size: 10px;
		font-weight: 500;
		border: 1px solid var(--color-border);
		transition: transform 0.15s ease;

		&.facedown {
			background: linear-gradient(
				135deg,
				color-mix(in srgb, var(--color-primary) 6%, var(--color-surface)),
				color-mix(in srgb, var(--color-primary) 15%, var(--color-surface))
			);
			border-color: color-mix(
				in srgb,
				var(--color-primary) 20%,
				var(--color-border)
			);
			box-shadow: 0 2px 8px
				color-mix(in srgb, var(--color-primary) 10%, transparent);
		}

		&.own {
			border-color: color-mix(
				in srgb,
				var(--color-primary) 55%,
				var(--color-border)
			);
			background: color-mix(
				in srgb,
				var(--color-primary) 12%,
				var(--color-surface)
			);
			box-shadow: 0 2px 10px
				color-mix(in srgb, var(--color-primary) 18%, transparent);
			cursor: default;

			&:hover {
				transform: translateY(-2px);
				box-shadow: 0 4px 14px
					color-mix(in srgb, var(--color-primary) 25%, transparent);
			}
		}

		&.revealed {
			opacity: 0.4;
			background: var(--color-surface);
			filter: grayscale(0.5);
			border-style: dashed;
		}
	}

	.card-back-pattern {
		font-size: 18px;
		color: var(--color-primary);
		opacity: 0.4;
	}

	.card-icon {
		font-size: 16px;
	}

	.card-label {
		font-size: 8px;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--color-text-secondary);
	}

	/* Role colors */
	.role-courtesan .card-icon {
		color: #f43f8c;
		filter: drop-shadow(0 0 3px rgba(244, 63, 140, 0.4));
	}
	.role-mercenary .card-icon {
		color: #f43f5e;
		filter: drop-shadow(0 0 3px rgba(244, 63, 94, 0.4));
	}
	.role-baron .card-icon {
		color: #b44dff;
		filter: drop-shadow(0 0 3px rgba(180, 77, 255, 0.4));
	}
	.role-consul .card-icon {
		color: #4d94ff;
		filter: drop-shadow(0 0 3px rgba(77, 148, 255, 0.4));
	}
	.role-admiral .card-icon {
		color: #0ed3cf;
		filter: drop-shadow(0 0 3px rgba(14, 211, 207, 0.4));
	}
</style>
