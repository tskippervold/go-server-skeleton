export { default as CustomerForm } from '../../components/CustomerForm.vue'
export { default as PaymentSelector } from '../../components/PaymentSelector.vue'
export { default as SelectedCustomer } from '../../components/SelectedCustomer.vue'

export const LazyCustomerForm = import('../../components/CustomerForm.vue' /* webpackChunkName: "components/CustomerForm" */).then(c => c.default || c)
export const LazyPaymentSelector = import('../../components/PaymentSelector.vue' /* webpackChunkName: "components/PaymentSelector" */).then(c => c.default || c)
export const LazySelectedCustomer = import('../../components/SelectedCustomer.vue' /* webpackChunkName: "components/SelectedCustomer" */).then(c => c.default || c)
