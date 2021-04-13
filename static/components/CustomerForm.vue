<template>
    <form id="customer-form">

        <label for="customer-type" class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2">
            Customer type
        </label>
        <div class="inline-block relative w-full mb-6">
            <select v-model="customer.customerType" id="customer-type" class="block appearance-none w-full bg-white border border-gray-400 hover:border-gray-500 px-4 py-2 pr-8 rounded shadow leading-tight focus:outline-none focus:shadow-outline">
                <option v-bind:value="'person'" selected>Private person</option>
                <option v-bind:value="'business'">Business</option>
            </select>
            <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                <svg class="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/></svg>
            </div>
            <p v-if="this.errors.customer" class="text-red-500 text-xs italic">{{ this.errors.customer }}</p>
        </div>

        <!---->

        <div class="flex flex-wrap -mx-3 mb-6">
            <div class="w-full md:w-1/2 px-3 mb-6 md:mb-0">
                <label class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" for="customer-first-name">
                    First Name
                </label>
                <input v-model="customer.firstName" class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="customer-first-name" type="text">
                <p v-if="this.errors.firstName" class="text-red-500 text-xs italic">{{ this.errors.firstName }}</p>
            </div>
            <div class="w-full md:w-1/2 px-3">
                <label class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" for="customer-last-name">
                    Last Name
                </label>
                <input v-model="customer.lastName" class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="customer-last-name" type="text">
                <p v-if="this.errors.lastName" class="text-red-500 text-xs italic">{{ this.errors.lastName }}</p>
            </div>
        </div>

        <div class="flex flex-wrap -mx-3 mb-6">
            <div class="w-full px-3">
                <label class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" for="customer-email">
                    Email
                </label>
                <input v-model="customer.email" class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="customer-email" type="email">
                <p v-if="this.errors.email" class="text-red-500 text-xs italic">{{ this.errors.email }}</p>
            </div>
        </div>

        <div class="w-full flex justify-center">
            <button
                v-for="method in this.paymentMethods"
                :key="method.id"
                class="text-base rounded py-2 px-6 mx-2 text-white shadow-lg hover:shadow-none self-end"
                style="background-color:#327dff"
                :disabled="method.selected"
                v-on:click.prevent="submitForm(method)">
                {{ method.title }}
            </button>
        </div>
    </form>
</template>

<script>
export default {
    data() {
        return {
            customer: {
                customerType: 'person',
                firstName: null,
                lastName: null,
                email: null
            },
            errors: {},
            paymentMethods: this.$store.state.paymentMethods
        }
    },
    methods: {
        submitForm(selectedPaymentMethod) {
            // Reset `error` object
            this.errors = {}

            if (!this.customer.customerType) {
                this.errors['customer'] = 'Please select a customer type'
            }

            if (!this.customer.firstName) {
                this.errors['firstName'] = 'Please enter your first name'
            }

            if (!this.customer.lastName) {
                this.errors['lastName'] = 'Please enter your last name'
            }

            if (!this.customer.email) {
                this.errors['email'] = 'Please enter a valid email'
            }

            if (Object.keys(this.errors).length != 0) {
                return
            }

            this.$store.commit('selectPaymentMethod', selectedPaymentMethod)
            this.$store.commit('setCustomer', this.customer)
        }
    }
}
</script>