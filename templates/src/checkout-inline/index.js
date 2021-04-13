/*import { InlineCheckout, Event } from "@bambora/checkout-sdk-web";

document.addEventListener('DOMContentLoaded', createCheckoutForm)

async function createSession() {
    const response = await fetch('/wc-api/wc_zo_gateway_card?action=create_session')
    return await response.json()
}

async function createCheckoutForm() {
    // https://v1.checkout.bambora.com
    try {
        const session = await createSession()
        console.log('SESSION CREATED', session)
    
        const token = session.token
        const mountContainer = document.getElementById('zo-cc-container')
        if (!mountContainer) {
            throw new Error('Could not find checkout container to mount.')
        }
    
        const checkout = new InlineCheckout(token, {
            eventHandlerMap: {
                [Event.Authorize]: [(e) => console.log(e)]
            }
        })

        await checkout.mount(mountContainer)
    } catch (error) {
        console.error(error)
    }
}*/

