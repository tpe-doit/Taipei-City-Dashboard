<script setup>
// Mapbox Related Imports
import mapboxGl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css'
import mapStyle from '../../assets/configs/mapbox/mapStyle.js'
import { MapPositions } from '../../assets/configs/mapbox/mapType.js'
import { MapObjectConfig, MapFogStyle, BuildingsIn3D } from '../../assets/configs/mapbox/mapConfig.js'

import { onMounted, ref } from 'vue';
import { useMapStore } from '../../store/mapStore';

const MAPBOXTOKEN = import.meta.env.VITE_MAPBOXTOKEN

const mapStore = useMapStore()

const mapBoxObject = ref(null)
// function initializeMapBox() {
//     mapboxGl.accessToken = MAPBOXTOKEN
//     mapBoxObject.value = new mapboxGl.Map({
//         ...MapObjectConfig,
//         style: mapStyle
//     })
//     mapBoxObject.value.addControl(new mapboxGl.NavigationControl())
//     mapBoxObject.value.doubleClickZoom.disable()
//     mapBoxObject.value.on('style.load', () => {

//         mapStore.initializeBasicLayers(mapBoxObject.value)
//         mapStore.initialize3DLayers(mapBoxObject.value)
//     })
// }
onMounted(() => {
    mapStore.initializeMapBox()
})
</script>

<template>
    <div class="mapcontainer">
        <div id="mapboxBox" />
    </div>
</template>

<style scoped lang="scss">
.mapcontainer {
    width: 100%;
    height: calc(100% - var(--font-m));
    background-color: var(--color-component-background);
    flex: 1;
    border-radius: 5px;
}

#mapboxBox {
    width: 100%;
    height: 100%;
}
</style>