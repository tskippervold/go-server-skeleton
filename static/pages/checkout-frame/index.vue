<template>
    <div class="max-w-lg mx-auto my-8">
        <section id="customer" class="rounded-lg p-4 mb-4" v-bind:class="hasCustomer ? '' : 'shadow'">
            <SelectedCustomer v-if="hasCustomer" v-bind:customer="this.customer" />
            <CustomerForm v-else />
        </section>

        <section v-if="hasCustomer" id="payment" class="shadow">
            <PaymentSelector />
        </section>
    </div>
</template>

<script>
export default {
    name: 'checkout-frame',
    mounted () {
        const bamboraToken = this.$route.query.bambora_token
        if (bamboraToken) {
            this.$store.commit('addBamboraPaymentMethod', bamboraToken)
        }
    },
    computed: {
        customer() {
            return this.$store.state.customer
        },
        hasCustomer() {
            return this.$store.state.customer != null
        }
    }
}
</script>
