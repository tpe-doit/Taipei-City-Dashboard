<script setup>
import { defineProps, computed } from "vue";

const props = defineProps({
	value: {
		type: Number,
		required: true,
	},
});

// Compute the color based on the value
const indicatorColor = computed(() => {
	if (props.value > 0 && props.value <= 33) {
		return "#62C554"; // green
	} else if (props.value >= 34 && props.value <= 66) {
		return "#F4BF4F"; // yellow
	} else if (props.value >= 67 && props.value <= 100) {
		return "#ED6A5D"; // red
	} else {
		return "var(--color-component-background)";
	}
});
</script>

<template>
	<div
		class="security-indicator"
		:class="{ 'is-inactive': value === 0 }"
		:style="{ backgroundColor: indicatorColor }"
	>
		<div>交通風險</div>
		<div class="security-indicator-value">
			{{ value }}
		</div>
	</div>
</template>

<style scoped lang="scss">
.security-indicator {
	position: absolute;
	top: var(--font-s);
	left: var(--font-s);
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 5px;
	width: 80px;
	height: 80px;
	color: white;
	font-size: var(--font-ms);
	text-align: center;
	padding: 8px;
	border-radius: 8px;
	z-index: 2;
	transition: background-color 0.2s, opacity 0.2s;
	pointer-events: none;

	&.is-inactive {
		opacity: 0;
		pointer-events: none;
	}

	&-value {
		flex: 1;
		display: flex;
		justify-content: center;
		align-items: center;
		width: 100%;
		font-size: var(--font-xl);
		color: var(--color-background);
		background-color: #fff;
		border-radius: 8px;
	}
}
</style>
