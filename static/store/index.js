export const state = () => ({
    customer: null,
    paymentMethods: [
        {
            selected: false,
            id: 'invoice-payment',
            title: 'Pay later',
            token: null
        }
    ]
})

export const mutations = {
    setCustomer(state, value) {
        state.customer = value
    },
    addBamboraPaymentMethod(state, token) {
        state.paymentMethods.push({
            selected: false,
            id: 'card-payment',
            title: 'Pay with card',
            token: token
        })
    },
    selectPaymentMethod(state, method) {
        state.paymentMethods.forEach(element => {
            if (element === method) {
                return
            }

            element.selected = false
        })

        if (method) {
            method.selected = true
        }
    }
}