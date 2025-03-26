window.addEventListener('load', function () {
  // Initialize Stripe
  const stripe = Stripe(window.env.STRIPE_PUBLISHABLE_KEY);
  let elements;
  
  // Cache DOM elements
  const form = document.getElementById('form');
  const sessionType = document.getElementById('session-type');
  const startTime = document.getElementById('start-time');
  const endTime = document.getElementById('end-time');
  const cymbals = document.getElementById('cymbals');
  const stripeSubmitBtn = document.getElementById('stripe-submit');
  const paymentContainer = document.getElementById('stripe-form');
  const formContainer = document.getElementById('form-container');
  const successMessage = document.getElementById('success');
  const pageHeader = document.getElementById('page-header');

  // Add event listeners for input changes
  sessionType.addEventListener('change', updatePrice);
  startTime.addEventListener('change', updatePrice);
  endTime.addEventListener('change', updatePrice);
  cymbals.addEventListener('change', updatePrice);
  
  // Initial price update
  updatePrice();
  
  function updatePrice() {
    fetchPrice().then(function (price) {
      stripeSubmitBtn.textContent = `Pay: £${price.toFixed(2)}`;
    });
  }
  
  async function createBooking() {
    const bookingData = {
      type: sessionType.value,
      name: document.getElementById('name').value,
      email: document.getElementById('email').value,
      phone: document.getElementById('phone').value,
      room: document.getElementById('room').value,
      date: document.getElementById('date-input').value,
      start_time: startTime.value,
      end_time: endTime.value,
    };
    
    const bookingResponse = await fetch('/bookings', {
      method: 'POST',
      body: JSON.stringify(bookingData),
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    if (!bookingResponse.ok) {
      throw new Error(`Can't book ${bookingData.room} at that time!`);
    }
    
    return bookingResponse.json();
  }
  
  async function getPaymentIntent(bookingId) {
    const intentResponse = await fetch('/stripe/payment-intents', {
      method: 'POST',
      body: JSON.stringify({
        booking_id: bookingId
      }),
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    return intentResponse.json();
  }
  
  async function confirmBooking(bookingId, paymentResult) {
    return fetch(`/bookings/${bookingId}/confirm`, {
      method: 'POST',
      body: JSON.stringify({
        payment_id: paymentResult.paymentIntent.id,
        payment_status: paymentResult.paymentIntent.status,
        payment_amount: paymentResult.paymentIntent.amount,
        payment_method: paymentResult.paymentIntent.payment_method
      }),
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }
  
  async function deleteBooking(bookingId) {
    return fetch(`/bookings/${bookingId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }
  
  function showError(message) {
    alert(message);
    console.error(message);
  }
  
  function setupStripeForm(clientSecret) {
    elements = stripe.elements({
      clientSecret: clientSecret,
      appearance: {
        theme: 'stripe'
      }
    });
    
    const paymentElement = elements.create('payment');
    paymentElement.mount('#payment-gateway');
  }
  
  async function handlePaymentSuccess(bookingId, paymentResult) {
    try {
      const confirmResponse = await confirmBooking(bookingId, paymentResult);
      
      if (!confirmResponse.ok) {
        throw new Error('Failed to confirm booking on the server');
      }
      
      // Update UI for success
      
      formContainer.style.display = 'none';
      paymentContainer.style.display = 'none';
      successMessage.style.display = 'block';
      pageHeader.textContent = 'SUCCESS! 🎉'
    } catch (error) {
      showError(`Payment was successful, but there was an error confirming your booking: ${error.message}`);
    }
  }
  
  async function handlePaymentFailure(bookingId, error) {
    try {
      // Delete the booking
      await deleteBooking(bookingId);
      showError(`Payment failed: ${error.message}`);
    } catch (deleteError) {
      showError(`Payment failed: ${error.message}. Additionally, there was an error cleaning up your booking: ${deleteError.message}`);
    }
  }
  
  // Handle form submission
  form.addEventListener('submit', async function (e) {
    e.preventDefault();
    
    try {
      // Show the payment form
      paymentContainer.style.display = 'block';
      
      // Create held booking
      const booking = await createBooking();
      const bookingId = booking.id;
      
      // Get payment intent from server
      const intent = await getPaymentIntent(bookingId);
      
      // Initialize Stripe Elements with client secret
      setupStripeForm(intent.client_secret);
      
      // Handle payment submission
      stripeSubmitBtn.onclick = async function (e) {
        e.preventDefault();
        
        const result = await stripe.confirmPayment({
          elements,
          confirmParams: {
            return_url: window.location.origin + '/booking-confirmation',
          },
          redirect: 'if_required',
        });
        
        if (result.error) {
          await handlePaymentFailure(bookingId, result.error);
        } else {
          await handlePaymentSuccess(bookingId, result);
        }
      };
      
    } catch (error) {
      showError(error.message);
    }
  });
});
