<script setup>
import { computed, ref } from "vue";

const props = defineProps(["chart_config", "activeChart", "series"]);

const thisMonth = ref("06");
const thisYear = ref("2023");
// const q1 = ref(null);
// const q3 = ref(null);

// const percentile = (arr, p) => {
// 	if (arr.length === 0) return 0;
// 	if (typeof p !== "number") throw new TypeError("p must be a number");
// 	if (p <= 0) return arr[0];
// 	if (p >= 1) return arr[arr.length - 1];

// 	const index = (arr.length - 1) * p,
// 		lower = Math.floor(index),
// 		upper = lower + 1,
// 		weight = index % 1;

// 	if (upper >= arr.length) return arr[lower];
// 	return arr[lower] * (1 - weight) + arr[upper] * weight;
// };

const calculatePercentile = (arr, value) => {
	let count = 0;
	for (let i = 0; i < arr.length; i++) {
		if (arr[i] <= value) {
			count++;
		}
	}
	return (count / arr.length) * 100;
};

const parsedSeries = computed(() => {
	let history = {};
	const toParse = props.series[0].data; // 假设这是从JSON文件中提取的数据
	for (let i = 0; i < toParse.length; i++) {
		let yearMonth = String(toParse[i]["年月"]);
		let year = yearMonth.slice(0, 4); // 获取年份
		if (yearMonth.slice(-2) === thisMonth.value) {
			// 检查是否为6月份的数据
			if (!history[year]) {
				history[year] = 0; // 如果该年份还没有记录，则初始化为0
			}
			history[year] += toParse[i].total; // 将该年6月的降水总量累加
		}
	}

	let historicalData = [];
	for (let year in history) {
		// if (year == thisYear.value) {
		// 	// 如果是今年的数据，则不显示
		// 	continue;
		// }
		historicalData.push(history[year]);
	}

	// historicalData.sort((a, b) => a - b);
	// q1.value = percentile(historicalData, 0.25);
	// q3.value = percentile(historicalData, 0.75);

	let rainIndex = Math.ceil(
		calculatePercentile(historicalData, history[thisYear.value])
	);

	return [rainIndex];
});

const level = computed(() => {
	if (parsedSeries.value <= 50) {
		return "低";
	} else if (parsedSeries.value <= 80) {
		return "中";
	} else {
		return "高";
	}
});

const chartOptions = ref({
	chart: {
		type: "radialBar",
		sparkline: {
			enabled: true,
		},
	},
	colors: [props.chart_config.alert_color[level.value]],
	plotOptions: {
		radialBar: {
			startAngle: -90,
			endAngle: 90,
			dataLabels: {
				show: true,
				name: {
					show: false,
				},
				value: {
					offsetY: -3,
					fontSize: "22px",
					color: "#E1E1E1",
				},
			},
		},
	},
	labels: ["Average Results"],
	legend: {
		show: true,
		position: "top",
		horizontalAlign: "center",
		labels: {
			useSeriesColors: false,
		},
		customLegendItems: ["低", "中", "高"],
		markers: {
			size: 0,
			fillColors: ["#31BD00", "#FF9110", "#E93838"],
		},
		formatter: function (seriesName) {
			return (
				seriesName + ":  " + props.chart_config.alert_range[seriesName]
			);
		},
	},
});
</script>

<template>
	<div v-if="activeChart === 'SpeedChart'" class="speedchart">
		<apexchart
			width="100%"
			type="radialBar"
			:options="chartOptions"
			:series="parsedSeries"
		></apexchart>
		<div class="speedchart-title">
			<h2 :style="{ color: props.chart_config.alert_color[level] }">
				風險程度：{{ level }}
			</h2>
		</div>
	</div>
</template>

<style scoped lang="scss">
.speedchart {
	&-title {
		display: flex;
		align-items: center;
		justify-content: center;
		margin-top: -15px;
		h2 {
			font-size: 1.5rem;
		}
	}
}
</style>
