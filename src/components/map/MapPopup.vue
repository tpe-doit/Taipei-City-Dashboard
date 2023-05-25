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
    animation: easein 0.2s linear;
    min-width: 200px;
    max-width: 400px;
    width: max-content;
}

.mapboxgl-popup-content {
    background-color: var(--color-component-background) !important;
    border: solid 1px var(--color-border);
    box-shadow: 0px 0px 10px rgb(35, 35, 35) !important;
    border-radius: 5px !important;
    padding: 10px !important;

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
    font-size: 1.2rem;
    color: var(--color-complement-text);
    right: 15px !important;
    top: 10px !important;
    line-height: 1rem;
}

.mappopup {
    h2 {
        color: var(--color-complement-text);
        margin-bottom: 0.5rem;
    }

    max-height: 200px;
    overflow-y: scroll;

    &-content {
        display: flex;

        width: fit-content;

        h3 {
            min-width: 80px;
        }

        p {
            text-align: justify;
        }
    }
}
</style>