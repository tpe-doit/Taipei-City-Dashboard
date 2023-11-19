<script setup>
import { useMapStore } from "../../store/mapStore";

const props = defineProps(["chart_config", "series", "map_config"]);
const mapStore = useMapStore();

/** @type {import("apexcharts").ApexOptions} */
const chartOptions = {
	chart: {
		toolbar: {
			tools: {
				download: false,
				pan: false,
				reset: "<p>" + "重置" + "</p>",
				zoomin: false,
				zoomout: false,
			},
		},
	},
	dataLabels: {
		enabled: false,
	},
	grid: {
		show: false,
	},
	legend: {
		show: props.series.length > 1 ? true : false,
	},
	markers: {
		hover: {
			size: 5,
		},
		size: 3,
		strokeWidth: 0,
	},
	tooltip: {
		custom: function ({ series, seriesIndex, dataPointIndex, w }) {
			// The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css
			return (
				'<div class="chart-tooltip">' +
				"<h6>" +
				`${w.config.series[seriesIndex].data[dataPointIndex].x}` +
				"</h6>" +
				"<span>" +
				series[seriesIndex][dataPointIndex] +
				` ${props.chart_config.unit}` +
				"</span>" +
				"</div>"
			);
		},
	},
	xaxis: {
		axisBorder: {
			color: "#555",
			height: "0.8",
		},
		axisTicks: {
			color: "#555",
		},
		crosshairs: {
			show: false,
		},
		tooltip: {
			enabled: false,
		},
	},
	stroke: {
		width: 1,
		curve: "smooth",
	},
	fill: {
		type: "gradient",
		gradient: {
			shade: "dark",
			gradientToColors: ["#99aaee"],
			shadeIntensity: 1,
			type: "horizontal",
			opacityFrom: 1,
			opacityTo: 1,
			stops: [0, 100],
		},
	},
};

function handleDataSelection(_, { xaxis }) {
	const { max, min } = xaxis;
	const id = "youbike_crowded-circle";
	if (typeof max === "number" && typeof min === "number") {
		// console.log({ max, min },[
		// 	"all",
		// 	[">=", ["get", "available_rent_prob_median"], min / 20],
		// 	["<=", ["get", "available_rent_prob_median"], max / 20],
		// ])
		mapStore.map.setFilter(id, [
			"all",
			[">=", ["get", "available_rent_prob_median"], min / 20],
			["<=", ["get", "available_rent_prob_median"], max / 20],
		]);
	} else {
		mapStore.clearLayerFilter(id);
	}
}
</script>

<template>
	<div>
		<apexchart
			width="100%"
			height="280px"
			type="area"
			:options="chartOptions"
			:series="series"
			@zoomed="handleDataSelection"
		></apexchart>
	</div>
</template>
