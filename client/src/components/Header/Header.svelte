<script lang="ts">
	import logo from '../../assets/images/rebellion.png'
	import ProfileMenu from './ProfileMenu.svelte'
	import MobileNav from './MobileNav.svelte'
</script>

<header class="header">
	<div class="header-glow"></div>

	<a href="#/" class="logo-wrap">
		<img src={logo} alt="Rebellion" class="logo" />
	</a>

	<nav class="nav-center">
		<button class="nav-link">
			<span class="nav-link-text">About Us</span>
		</button>
		<button class="nav-link">
			<span class="nav-link-text">How to Play</span>
		</button>
	</nav>

	<div class="nav-right">
		<ProfileMenu />
	</div>

	<MobileNav />

	<div class="header-border"></div>
</header>

<style lang="scss">
	.header {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 16px 0 0;
		height: 72px;
		background: color-mix(in srgb, var(--color-bg) 80%, transparent);
		backdrop-filter: blur(16px);
		-webkit-backdrop-filter: blur(16px);
		z-index: 100;
	}

	.header-glow {
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 60%;
		height: 1px;
		background: radial-gradient(
			ellipse at center,
			var(--color-primary) 0%,
			transparent 70%
		);
		opacity: 0.5;
		pointer-events: none;
	}

	.header-border {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		height: 1px;
		background: linear-gradient(
			90deg,
			transparent 0%,
			var(--color-border) 20%,
			var(--color-primary) 50%,
			var(--color-border) 80%,
			transparent 100%
		);
		opacity: 0.6;
	}

	.logo-wrap {
		display: flex;
		align-items: center;
		text-decoration: none;
		transition: filter 0.3s ease;

		&:hover {
			filter: drop-shadow(
				0 0 8px color-mix(in srgb, var(--color-primary) 40%, transparent)
			);
		}
	}

	.logo {
		height: 48px;
		width: auto;

		:global([data-theme='dark']) & {
			filter: invert(1) hue-rotate(180deg);
		}
	}

	.nav-center {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.nav-link {
		position: relative;
		background: none;
		border: none;
		color: var(--color-text-secondary);
		font-family: inherit;
		font-size: 12px;
		font-weight: 500;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		cursor: pointer;
		padding: 8px 16px;
		border-radius: 6px;
		transition:
			color 0.2s ease,
			background 0.2s ease;

		&::after {
			content: '';
			position: absolute;
			bottom: 2px;
			left: 50%;
			transform: translateX(-50%) scaleX(0);
			width: 60%;
			height: 2px;
			background: var(--color-primary);
			border-radius: 1px;
			transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
			box-shadow: 0 0 8px
				color-mix(in srgb, var(--color-primary) 50%, transparent);
		}

		&:hover {
			color: var(--color-text);
			background: color-mix(in srgb, var(--color-primary) 6%, transparent);

			&::after {
				transform: translateX(-50%) scaleX(1);
			}
		}
	}

	.nav-right {
		display: flex;
		align-items: center;
	}

	@media (max-width: 550px) {
		.nav-center,
		.nav-right {
			display: none;
		}
	}
</style>
