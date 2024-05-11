<script setup>
import { defineProps, defineEmits, ref } from "vue";

const props = defineProps({
	placeholder: {
		type: String,
		default: "Search...",
	},
});

const emit = defineEmits(["search"]);

// Define reactive variables
const searchQuery = ref("");

// Function to handle search
function handleSearch() {
	emit("search", searchQuery.value);
}
</script>

<template>
	<div class="searchbar">
		<div class="searchbar-input">
			<input
				type="text"
				v-model="searchQuery"
				:placeholder="props.placeholder"
				@keydown.enter="handleSearch"
			/>
			<span
				v-if="searchQuery !== ''"
				@click="
					() => {
						searchQuery = '';
					}
				"
				>cancel</span
			>
		</div>
		<button @click="handleSearch">搜尋</button>
	</div>
</template>

<style scoped lang="scss">
.searchbar {
	display: flex;
	column-gap: 0.5rem;

	&-input {
		flex: 1;
		position: relative;
		display: flex;

		input {
			width: 100%;
		}

		span {
			position: absolute;
			right: 0;
			top: 0.3rem;
			margin-right: 4px;
			color: var(--color-complement-text);
			font-family: var(--font-icon);
			font-size: var(--font-m);
			cursor: pointer;
			transition: color 0.2s;

			&:hover {
				color: var(--color-highlight);
			}
		}
	}

	button {
		flex: none;
		display: flex;
		align-items: center;
		justify-self: baseline;
		padding: 0px 4px;
		border-radius: 5px;
		background-color: var(--color-highlight);
		font-size: var(--font-ms);
		transition: opacity 0.2s;

		&:hover {
			opacity: 0.8;
		}
	}
}
</style>
