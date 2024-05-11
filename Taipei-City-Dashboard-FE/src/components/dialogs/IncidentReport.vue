<!-- Developed by Taipei Urban Intelligence Center 2023-2024-->

<script setup>
import { ref, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useDialogStore } from "../../store/dialogStore";
import { useAuthStore } from "../../store/authStore";
import http from "../../router/axios";

import DialogContainer from "./DialogContainer.vue";

const dialogStore = useDialogStore();
const authStore = useAuthStore();

const location = ref(null);
const errorMessage = ref(null);
const incidentType = ref("");
const incidentDesc = ref("");
const incidentDis = ref("");

const { editUser } = storeToRefs(authStore);

const typeOptions = [
	{ label: "火災 Fire", value: "fire" },
	{ label: "淹水 Flood", value: "flood" },
	{ label: "道路 Road", value: "road" },
	{ label: "建物 Building", value: "building" },
	{ label: "其他 Others", value: "other" },
	// Add more options as needed
];

const disOptions = [
	{ label: "500公尺內", value: 0.5 },
	{ label: "500公尺~2公里", value: 2 },
	{ label: "2公里~5公里", value: 5 },
	{ label: "大於5公里", value: 10 },
	// Add more options as needed
];

function handleClose() {
	dialogStore.hideAllDialogs();
}

const getCurrentLocation = () => {
	if (navigator.geolocation) {
		navigator.geolocation.getCurrentPosition(
			(position) => {
				location.value = {
					latitude: position.coords.latitude,
					longitude: position.coords.longitude,
				};
			},
			(error) => {
				errorMessage.value = error.message;
			}
		);
		console.log(location);
	} else {
		errorMessage.value = "Geolocation is not supported by this browser.";
	}
};

async function handleSubmit() {
	let payload = {
		inctype: incidentType.value,
		description: incidentDesc.value,
		distance: incidentDis.value,
		latitude: location.value.latitude,
		longitude: location.value.longitude,
		reportTime: new Date().getTime(),
	};
	if (
		incidentType.value == "" ||
		incidentDesc.value == "" ||
		incidentDis.value == ""
	) {
		dialogStore.showNotification("fail", "Some properties are empty");
		return;
	}
	console.log(payload);
	const response = await http.post("/incident/", payload);
	console.log(response);
	incidentType.value = "";
	incidentDesc.value = "";
	incidentDis.value = "";
	dialogStore.showNotification("success", "OK!!!");
	dialogStore.hideAllDialogs();
}
onMounted(() => {
	getCurrentLocation();
});
</script>

<template>
	<DialogContainer :dialog="`incidentReport`" @on-close="handleClose">
		<div class="incidentreport">
			<h2>事件通報</h2>
			<label> 事件類型 </label>
			<select v-model="incidentType">
				<option
					v-for="(option, index) in typeOptions"
					:key="index"
					:value="option.value"
				>
					{{ option.label }}
				</option>
			</select>

			<label> 事件描述 </label>
			<input
				type="text"
				placeholder="(請概述事件過程)"
				v-model="incidentDesc"
			/>
			<label> 事件發生位置 </label>
			<select v-model="incidentDis">
				<option
					v-for="(option, index) in disOptions"
					:key="index"
					:value="option.value"
				>
					{{ option.label }}
				</option>
			</select>
			<label> 通報位置 </label>
			<!-- <input :value="parseTime(editUser.login_at)" disabled /> -->
			<input
				:value="location.latitude + `, ` + location.longitude"
				disabled
			/>
			<label> 通報時間 </label>
			<input :value="new Date().toLocaleString()" disabled />
			<div class="incidentreport-control">
				<button
					class="incidentreport-control-confirm"
					@click="handleSubmit"
				>
					提交
				</button>
			</div>
		</div>
	</DialogContainer>
</template>

<style scoped lang="scss">
.incidentreport {
	width: 300px;
	display: flex;
	flex-direction: column;

	label {
		margin: 8px 0 4px;
		font-size: var(--font-s);
		color: var(--color-complement-text);
	}

	&-control {
		height: 27px;
		display: flex;
		justify-content: flex-end;
		margin-top: var(--font-ms);

		&-confirm {
			margin: 0 2px;
			padding: 4px 10px;
			border-radius: 5px;
			background-color: var(--color-highlight);
			transition: opacity 0.2s;

			&:hover {
				opacity: 0.8;
			}
		}
	}
}
</style>
