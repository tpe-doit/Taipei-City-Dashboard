<script setup>
import { ref } from "vue";
import { useMapStore } from "../../store/mapStore";

const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);
const mapStore = useMapStore();

const chartOptions = ref({
	chart: {
		type: "bar",
		height: 350,
	},
	plotOptions: {
		bar: {
			borderRadius: 0,
			horizontal: true,
			barHeight: "80%",
			isFunnel: true,
		},
	},
	dataLabels: {
		enabled: true,
		formatter: function (val, opt) {
			return opt.w.globals.labels[opt.dataPointIndex] + ":  " + val;
		},
		dropShadow: {
			enabled: true,
		},
	},
	title: {
		text: "Recruitment Funnel",
		align: "center",
	},
	xaxis: {
		categories: [
			"Sourced",
			"Screened",
			"Assessed",
			"HR Interview",
			"Technical",
			"Verify",
			"Offered",
			"Hired",
		],
	},
	legend: {
		show: false,
	},
});

const chartHeight = 350;

const selectedIndex = ref(null);

function handleDataSelection(e, chartContext, config) {
	if (!props.chart_config.map_filter) {
		return;
	}
	if (config.dataPointIndex !== selectedIndex.value) {
		mapStore.addLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`,
			props.chart_config.map_filter[0],
			props.chart_config.map_filter[1][config.dataPointIndex]
		);
		selectedIndex.value = config.dataPointIndex;
	} else {
		mapStore.clearLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`
		);
		selectedIndex.value = null;
	}
}
</script>

<template>
	<div v-if="activeChart === 'FunnelChart'">
		<apexchart
			width="100%"
			:height="chartHeight"
			type="bar"
			:options="chartOptions"
			:series="series"
		></apexchart>
	</div>
</template>
