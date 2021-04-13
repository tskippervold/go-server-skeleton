<template>
    <div id="selected-customer">
        <div class="flex items-center">
            <div class="flex-grow">
                <span class="rounded-full text-white uppercase px-2 py-1 text-xs font-bold mr-3" style="background-color:#327dff">
                    {{ customer.customerType == 'person' ? 'Private customer' : 'Business customer'}}
                </span>
                <p class="text-sm mt-3">
                    <span>{{ customer.firstName}} {{ customer.lastName }}</span></br>
                    <span>{{ customer.email }}</span>
                </p>
            </div>
            <div>
                <button class="text-base rounded py-2 text-gray-600" v-on:click="changeCustomer">Change</button>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        customer: {
            required: true
        }
    },
    methods: {
        changeCustomer(e) {
            e.preventDefault()

            this.$store.commit('selectPaymentMethod', null)
            this.$store.commit('setCustomer', null)
        }
    }
}
</script>