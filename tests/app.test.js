import { describe, expect, it } from "vitest";
import { mount } from "@vue/test-utils";
import { createTestingPinia } from "@pinia/testing";

import { useAuthStore } from "../src/store/authStore";
import { useContentStore } from "../src/store/contentStore";
// import { useDialogStore } from "../src/store/dialogStore";
import { useMapStore } from "../src/store/mapStore";
import router from "../src/router";

import App from "../src/App.vue";

describe("Initial Application State", async () => {
	mount(App, {
		global: {
			plugins: [createTestingPinia, router],
		},
	});

	const authStore = useAuthStore();
	const contentStore = useContentStore();
	const mapStore = useMapStore();

	router.push("/dashboard");
	await router.isReady();

	it("App.vue Calls Initiation Functions", () => {
		expect(authStore.checkIfMobile).toHaveBeenCalledTimes(1);
		expect(authStore.setUser).toHaveBeenCalledTimes(1);
	});

	it("contentStore Calls Initiation Functions", () => {
		expect(contentStore.setRouteParams).toHaveBeenCalled();
	});

	router.push("/mapview");
	await router.isReady();

	it("mapStore Calls Initiation Functions", () => {
		expect(mapStore.initializeMapBox).toHaveBeenCalledTimes(1);
	});
});
