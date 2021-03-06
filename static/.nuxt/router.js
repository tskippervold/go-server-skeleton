import Vue from 'vue'
import Router from 'vue-router'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _67f60864 = () => interopDefault(import('../pages/checkout-frame/index.vue' /* webpackChunkName: "pages/checkout-frame/index" */))
const _194c3fae = () => interopDefault(import('../pages/checkout-inline/index.vue' /* webpackChunkName: "pages/checkout-inline/index" */))
const _3555dbfd = () => interopDefault(import('../pages/index.vue' /* webpackChunkName: "pages/index" */))

// TODO: remove in Nuxt 3
const emptyFn = () => {}
const originalPush = Router.prototype.push
Router.prototype.push = function push (location, onComplete = emptyFn, onAbort) {
  return originalPush.call(this, location, onComplete, onAbort)
}

Vue.use(Router)

export const routerOptions = {
  mode: 'history',
  base: decodeURI('/'),
  linkActiveClass: 'nuxt-link-active',
  linkExactActiveClass: 'nuxt-link-exact-active',
  scrollBehavior,

  routes: [{
    path: "/checkout-frame",
    component: _67f60864,
    name: "checkout-frame"
  }, {
    path: "/checkout-inline",
    component: _194c3fae,
    name: "checkout-inline"
  }, {
    path: "/",
    component: _3555dbfd,
    name: "index"
  }],

  fallback: false
}

export function createRouter () {
  return new Router(routerOptions)
}
