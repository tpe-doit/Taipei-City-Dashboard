<!-- Developed by Taipei Urban Intelligence Center 2023 -->

<script setup>
import { ref, reactive } from "vue";

import MetroCarDensity from "../utilities/MetroCarDensity.vue";
import { lineInfo } from "../../assets/configs/apexcharts/taipeiMetroLines";
import { computed } from "vue";
import { isHTMLTag } from "@vue/shared";

const props = defineProps(["chart_config", "activeChart", "series"]);
const count = ref(1);
const parsedSeries = computed(() => {
	console.log({ count });
	const number = props.series[count.value].data * 100;
	console.log({ number });

	return fill1In2d(number);
});

function fill1In2d(n) {
	console.log(n);
	// Create a 10x10 2D array in JavaScript
	const twoDArray = [];

	// Loop to create rows
	for (let i = 0; i < 10; i++) {
		// Initialize an empty row
		const row = [];
		// Loop to create columns in each row
		for (let j = 0; j < 10; j++) {
			// Add an element to the row
			row.push(0);
		}
		// Add the row to the 2D array
		twoDArray.push(row);
	}

	for (let i = 0; i < n; i++) {
		twoDArray[Math.floor(i / 10)][i % 10] = 1;
	}
	return twoDArray;
}

function calculateCount(n) {
	console.log(n);
	count.value + n;
	console.log(count.value);
}

function increaseCount() {
	count.value++;
}
function decreaseCount(n) {
	count.value--;
}
function showOrHideButton(isShow) {
	return isShow ? "show" : "hide";
}

console.log({ parsedSeries });
</script>

<style scoped lang="scss">
.hundredchart {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	gap: 10px;
	.chartcontainer {
		display: flex;
		flex-direction: row;
		gap: 10px;
		button {
			img {
				width: 30px;
				height: 30px;
			}
		}
	}
	.chart {
		.row {
			display: flex;
			flex-direction: column;
			justify-content: center;
			display: flex;
			flex-direction: row;
			justify-content: center;
			gap: 5px;
			.img {
				width: 15px;
				height: 15px;
			}
		}
	}
	.hide {
		visibility: hidden;
	}
	.show {
		visibility: visible;
	}
}
</style>

<template>
	<div v-if="activeChart === 'HundredChart'" class="hundredchart">
		<div>
			每 {{ props.series[count - 1].data * 100 }}
			{{ props.series[count - 1].unit
			}}{{ props.series[count - 1].name }} 就有
			{{ series[count].data * 100 }} {{ props.series[count].unit
			}}{{ props.series[count].name }}
		</div>
		<div class="chartcontainer">
			<button @click="decreaseCount" :class="showOrHideButton(count > 1)">
				<img src="../../assets/images/hundredicon/arrowLeft.svg" />
			</button>
			<div class="chart">
				<div
					v-for="(row, i) in parsedSeries"
					class="row"
					:key="`${row}-${i}`"
				>
					<div v-for="(item, j) in row" :key="`item-${j}`">
						<img
							v-if="item === 0"
							src="../../assets/images/hundredicon/human.svg"
							:alt="`numerator-${i}-${j}-${item}`"
							class="img"
						/>
						<img
							v-else
							src="../../assets/images/hundredicon/girls.svg"
							:alt="`denominator-${i}-${j}`"
							class="img"
						/>
					</div>
				</div>
			</div>
			<button
				@click="increaseCount"
				:class="showOrHideButton(count < props.series.length - 1)"
			>
				<img src="../../assets/images/hundredicon/arrowRight.svg" />
			</button>
		</div>
	</div>
</template>
