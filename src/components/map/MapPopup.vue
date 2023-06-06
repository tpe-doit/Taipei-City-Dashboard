<!-- Cleaned -->

<!-- This component is mounted programmically by the mapstore. "mapConfig" and "popupContent" are passed in in the mapStore -->
<!-- !! Do not delete the script tag !! -->
<script setup>
</script>

<template>
    <div class="mappopup">
        <h2>{{ mapConfig.title }}</h2>
        <div>
            <div class="mappopup-content" v-for="item in mapConfig.property" :key="item.key">
                <h3>{{ item.name }}</h3>
                <p>{{ popupContent.properties[item.key] }}</p>
            </div>
        </div>

    </div>
</template>

<style lang="scss">
@keyframes easein {
    0% {
        opacity: 0
    }

    100% {
        opacity: 1
    }
}

.mapboxgl-popup {
    width: max-content;
    min-width: 200px;
    max-width: 400px;
    animation: easein 0.2s linear;
}

.mapboxgl-popup-content {
    padding: 10px !important;
    border: solid 1px var(--color-border);
    box-shadow: 0px 0px 10px rgb(35, 35, 35) !important;
    border-radius: 5px !important;
    background-color: var(--color-component-background) !important;
}

.mapboxgl-popup-anchor-bottom .mapboxgl-popup-tip,
.mapboxgl-popup-anchor-bottom-left .mapboxgl-popup-tip,
.mapboxgl-popup-anchor-bottom-right .mapboxgl-popup-tip {
    border-top-color: var(--color-border) !important;
}

.mapboxgl-popup-anchor-top .mapboxgl-popup-tip,
.mapboxgl-popup-anchor-top-left .mapboxgl-popup-tip,
.mapboxgl-popup-anchor-top-right .mapboxgl-popup-tip {
    border-bottom-color: var(--color-border) !important;
}

.mapboxgl-popup-anchor-left .mapboxgl-popup-tip {
    border-right-color: var(--color-border) !important;
}

.mapboxgl-popup-anchor-right .mapboxgl-popup-tip {
    border-left-color: var(--color-border) !important;
}

.mapboxgl-popup-close-button {
    right: 15px !important;
    top: 10px !important;
    color: var(--color-complement-text);
    font-size: 1.2rem;
    line-height: 1rem;
}

.mappopup {
    max-height: 200px;
    overflow-y: scroll;

    h2 {
        margin-bottom: 0.5rem;
        color: var(--color-complement-text);
    }

    &-content {
        width: fit-content;
        display: flex;

        h3 {
            min-width: 80px;
        }

        p {
            text-align: justify;
        }
    }
}
</style>