<template>
    <div id="payment-selector">
        <div class="w-full flex justify-center pt-6 mb-4">
            <button
                v-for="method in this.paymentMethods"
                :key="method.id"
                class="text-base rounded py-2 px-6 mx-2 hover:shadow-none self-end"
                v-bind:class="method.selected ? 'text-white shadow-lg' : 'text-black shadow'"
                v-bind:style="method.selected ? 'background-color:#327dff' : ''"
                :disabled="method.selected"
                v-on:click="selectPaymentMethod(method)">
                {{ method.title }}
            </button>
        </div>
        <div>
            <div class="px-10 pt-2 pb-4 text-grey-darkest text-sm">
                <div ref="external-mount" style="height: 500px"></div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            externalMountedObject: null,
            paymentMethods: this.$store.state.paymentMethods
        }
    },
    mounted() {
        this.paymentMethods.forEach(element => {
            if (element.selected) {
                this.didSelectPaymentMethod(element)
            }
        })
    },
    methods: {
        selectPaymentMethod(method) {
            this.$store.commit('selectPaymentMethod', method)
            this.didSelectPaymentMethod(method)
        },
        async didSelectPaymentMethod(method) {
            if (this.externalMountedObject) {
                this.externalMountedObject.destroy()
            }

            const checkout = new Bambora.InlineCheckout(method.token, {
                eventHandlerMap: {
                    "*": [function(eventType, payload) {
                        console.log(eventType, payload);
                    }]
                }
            })

            const mountTo = this.$refs['external-mount']
            await checkout.mount(mountTo)
            this.externalMountedObject = checkout
        }
    }
}
</script>
