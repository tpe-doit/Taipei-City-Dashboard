<!-- Developed by Taipei Urban Intelligence Center 2023 -->

<script setup>
import { onMounted, ref , computed } from "vue";
import { useMapStore } from "../../store/mapStore";
import { useContentStore } from "../../store/contentStore";

const { BASE_URL } = import.meta.env;

const props = defineProps(["chart_config", "series", "map_config","activeChart"]);
const mapStore = useMapStore();
const contentStore = useContentStore();

const selectedIndex = ref(null);
const sliderValue = ref(0);
const maxval = computed(() => {
	return props.chart_config.map_filter[1].length - 1;
});
function handleDataSelection(index) {
	if (!props.chart_config.map_filter) {
		return;
	}
	if (index !== selectedIndex.value) {
		mapStore.addLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`,
			props.chart_config.map_filter[0],
			props.chart_config.map_filter[1][index]
		);
		selectedIndex.value = index;
	} else {
		mapStore.clearLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`
		);
		selectedIndex.value = null;
	}
}


</script>

<template>
		<div class="maplegend-legend" v-if="activeChart === 'MapSlideBar'">
				<input style="width: 100%; height: 25px; padding: 0;" type="range" min="0" :max="maxval" v-model="sliderValue" @input="handleDataSelection(sliderValue)">
				<p style="text-align: center;font-size: large;"> {{ props.chart_config.map_filter[1][sliderValue] }}</p>
		</div>
</template>

<style scoped lang="scss">
.maplegend {
	flex-direction: column;
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-top: -1rem;
	overflow: visible;


	&-legend {
		margin-top:25px;
		width: 100%;
		display: flex;
		flex-direction: column;
		grid-template-columns: 1fr 1fr;
		column-gap: 0.5rem;
		row-gap: 0.5rem;
		overflow: visible;
		margin-top: 40%;

		&-item {
			display: flex;
			align-items: center;
			padding: 5px 10px 5px 5px;
			border: 1px solid transparent;
			border-radius: 5px;
			transition: box-shadow 0.2s;
			cursor: auto;

			div:first-child,
			img {
				width: 1rem;
				margin-right: 0.75rem;
			}

			h5 {
				color: var(--color-complement-text);
				font-size: 0.75rem;
				text-align: left;
			}

			h6 {
				font-size: 1rem;
				font-weight: 400;
				text-align: center;
			}

			p {
				font-size: 1rem;
				font-weight: 400;
				text-align: center;
			}
		}
	}

	&-filter {
		border: 1px solid var(--color-border);
		cursor: pointer;

		&:hover {
			box-shadow: 0px 0px 5px black;
		}
	}

	&-selected {
		box-shadow: 0px 0px 5px black;
	}


}
</style>
