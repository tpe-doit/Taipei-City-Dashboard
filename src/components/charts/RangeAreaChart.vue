<script setup>
import { ref, computed } from "vue";
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
		// offsetY: 15,
		stacked: false,
		toolbar: {
			show: false,
		},
	},
	colors: props.chart_config.color,
	dataLabels: {
		enabled: false,
		offsetX: 20,
		textAnchor: "start",
		style: {
			colors: ["#000000"],
		},
	},
	fill: {
		colors: ["#397AB7"],
	},
	grid: {
		show: false,
	},
	legend: {
		show: false,
	},
	plotOptions: {},
	stroke: {
		colors: ["#282a2c"],
		show: true,
		width: 0,
	},
	tooltip: {
		custom: function ({ series, seriesIndex, dataPointIndex, w }) {
			return (
				'<div class="chart-tooltip">' +
				"<h6>" +
				w.globals.labels[dataPointIndex] +
				"</h6>" +
				"<span>" +
				series[seriesIndex][dataPointIndex] +
				` ${props.chart_config.unit}` +
				"</span>" +
				"</div>"
			);
		},
		followCursor: true,
	},
	xaxis: {
		axisBorder: {
			show: false,
		},
		axisTicks: {
			show: false,
		},
		type: "category",
		labels: {
			style: {
				colors: ["#282a2c"],
			},
		},
	},
	// yaxis: {
	// 	labels: {
	// 		formatter: function (value) {
	// 			return value.length > 7 ? value.slice(0, 6) + "..." : value;
	// 		},
	// 	},
	// },
});

const processedSeries = computed(() => {
	const series = props.series[0].data;

	// console.log('係瑞斯', series)
	// series[0]['town'] = "你好！"
	// "town": "北投區",
	// "stationid": "'001'",
	// "name": "湖田國小",
	// "x坐標": 304329.0934,
	// "y坐標": 2784567.876,
	// "年月": 200007,
	// "total": 0.0

	const date_to_rain = {};
	for (const serie of series) {
		if (!date_to_rain[serie["年月"]]) {
			date_to_rain[serie["年月"]] = [];
		}
		date_to_rain[serie["年月"]].push(serie["total"]);
	}

	const calculateQuartiles = (data) => {
		const sortedData = data.slice().sort((a, b) => a - b);
		const medianIndex = Math.floor(sortedData.length / 2);

		const lowerHalf = sortedData.slice(0, medianIndex);
		const upperHalf = sortedData.slice(
			medianIndex + (sortedData.length % 2 === 0 ? 0 : 1)
		);

		const lowerQuartile = calculateMedian(lowerHalf);
		const median = calculateMedian(sortedData);
		const upperQuartile = calculateMedian(upperHalf);

		return {
			lowerQuartile,
			median,
			upperQuartile,
		};
	};

	const calculateMedian = (sortedData) => {
		const { length } = sortedData;
		if (length % 2 === 0) {
			return (sortedData[length / 2 - 1] + sortedData[length / 2]) / 2;
		} else {
			return sortedData[Math.floor(length / 2)];
		}
	};

	const date_to_stat = [];
	for (const key in date_to_rain) {
		// const (f, s, t) = calculateQuartiles(date_to_rain[key])
		date_to_stat.push({
			x: Number(key),
			y: [
				Math.min(...date_to_rain[key]) / 100,
				Math.max(...date_to_rain[key]) / 100,
			],
		});
	}

	return [
		{
			name: "",
			data: date_to_stat.slice(-12),
		},
	];
});
</script>

<template>
	<div v-if="activeChart === 'RangeAreaChart'">
		<apexchart
			width="100%"
			type="rangeArea"
			:options="chartOptions"
			:series="processedSeries"
		>
		</apexchart>
	</div>
</template>
