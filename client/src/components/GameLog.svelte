<script lang="ts">
	interface LogEntry {
		id: number
		text: string
		type: 'action' | 'challenge' | 'block' | 'elimination' | 'system'
	}

	let { entries }: { entries: LogEntry[] } = $props()
</script>

<div class="sidebar">
	<div class="game-log">
		<h2 class="panel-title">Game Log</h2>
		<div class="log-entries">
			{#each entries as entry}
				<div class="log-entry type-{entry.type}">
					<span class="log-text">{entry.text}</span>
				</div>
			{/each}
		</div>
	</div>
</div>

<style lang="scss">
	.sidebar {
		display: flex;
		flex-direction: column;
		border-left: 1px solid
			color-mix(in srgb, var(--color-primary) 15%, var(--color-border));
		background: linear-gradient(
			180deg,
			color-mix(in srgb, var(--color-primary) 4%, var(--color-surface)),
			var(--color-surface) 40%
		);
		overflow: hidden;
		height: 100%;
	}

	.panel-title {
		font-size: 11px;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: var(--color-primary);
		padding: 16px 16px 12px;
		text-shadow: 0 0 8px
			color-mix(in srgb, var(--color-primary) 30%, transparent);
		border-bottom: 1px solid
			color-mix(in srgb, var(--color-primary) 12%, var(--color-border));
	}

	.game-log {
		flex: 1;
		display: flex;
		flex-direction: column;
		min-height: 0;
	}

	.log-entries {
		flex: 1;
		overflow-y: auto;
		padding: 8px 12px 16px;
		display: flex;
		flex-direction: column;
		gap: 5px;
	}

	.log-entry {
		padding: 8px 12px;
		border-radius: 8px;
		font-size: 12px;
		line-height: 1.4;
		color: var(--color-text-secondary);
		border-left: 3px solid transparent;
		transition:
			background 0.15s ease,
			transform 0.1s ease;

		&:hover {
			transform: translateX(2px);
		}

		&.type-action {
			border-left-color: var(--color-primary);
			background: color-mix(in srgb, var(--color-primary) 8%, transparent);
			box-shadow: inset 0 0 12px
				color-mix(in srgb, var(--color-primary) 5%, transparent);
		}
		&.type-challenge {
			border-left-color: #ffb800;
			color: #ffb800;
			background: color-mix(in srgb, #ffb800 8%, transparent);
			box-shadow: inset 0 0 12px color-mix(in srgb, #ffb800 5%, transparent);
		}
		&.type-block {
			border-left-color: #0ed3cf;
			color: #0ed3cf;
			background: color-mix(in srgb, #0ed3cf 8%, transparent);
			box-shadow: inset 0 0 12px color-mix(in srgb, #0ed3cf 5%, transparent);
		}
		&.type-elimination {
			border-left-color: #f43f5e;
			color: #f43f5e;
			background: color-mix(in srgb, #f43f5e 8%, transparent);
			box-shadow: inset 0 0 12px color-mix(in srgb, #f43f5e 5%, transparent);
		}
		&.type-system {
			font-style: italic;
			color: var(--color-primary);
			opacity: 0.8;
			text-shadow: 0 0 6px
				color-mix(in srgb, var(--color-primary) 20%, transparent);
		}
	}

	.log-text {
		display: block;
		font-weight: 500;
	}

	@media (max-width: 900px) {
		.sidebar {
			border-left: none;
			border-top: 1px solid var(--color-border);
			max-height: 40vh;
		}
	}

	@media (max-width: 700px) {
		.sidebar {
			max-height: 25vh;
		}
	}
</style>
