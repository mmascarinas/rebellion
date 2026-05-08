<script lang="ts">
	let { on = false, onToggle }: { on: boolean; onToggle?: () => void } =
		$props()
</script>

<button
	class="toggle-switch"
	class:on
	onclick={onToggle}
	aria-label="Toggle switch"
>
	<span class="toggle-knob"></span>
</button>

<style lang="scss">
	.toggle-switch {
		position: relative;
		display: inline-block;
		width: 36px;
		height: 20px;
		background: var(--color-border);
		border: none;
		border-radius: 10px;
		cursor: pointer;
		transition: background 0.25s ease;
		flex-shrink: 0;
		padding: 0;

		&.on {
			background: var(--color-primary);
		}
	}

	.toggle-knob {
		position: absolute;
		display: block;
		top: 2px;
		left: 2px;
		width: 16px;
		height: 16px;
		background: var(--color-bg);
		border-radius: 50%;
		transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.15);
		pointer-events: none;

		.on & {
			transform: translateX(16px);
		}
	}
</style>
