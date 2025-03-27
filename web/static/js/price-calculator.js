async function fetchPrice() {
    let session_type = document.getElementById('session-type').value
    let date = document.getElementById('date-input').value
    let start_time = date + " " + document.getElementById('start-time').value
    let end_time = date + " " + document.getElementById('end-time').value
    let cymbals = document.getElementById('cymbals').checked ? 1 : 0

    if (document.getElementById('end-time').value != "" ) {
        try {
            const response = await fetch(`/price-calculator?startTime=${start_time}&endTime=${end_time}&type=${session_type}&cymbals=${cymbals}`);
            if (!response.ok) throw new Error('Failed to fetch availability');
            const data = await response.json();
            return data.price;
        } catch (error) {
            console.error(error);
            return [];
        }
    }

    return 0;
}

function updatePrice() {
    let price = document.getElementById('price')
    const stripeSubmitBtn = document.getElementById('stripe-submit');
    fetchPrice().then(data => {
        price.textContent = `£${(data).toFixed(2)}`
        if (stripeSubmitBtn) {
            stripeSubmitBtn.textContent = `Pay: £${(data).toFixed(2)}`; 
        }
    })
}

window.addEventListener('load', function () {
    populateSummary();
    let session_type = document.getElementById('session-type')
    let end_time = document.getElementById('end-time')
    let buttons = document.querySelectorAll('.toggle-button');
    let cymbals = document.getElementById('cymbals')

    // event listeners for input changes
    session_type.addEventListener('change', updatePrice)
    session_type.addEventListener('change', populateSummary)
    end_time.addEventListener('change', updatePrice)
    end_time.addEventListener('change', populateSummary)
    cymbals.addEventListener('change', updatePrice)
    cymbals.addEventListener('change', populateSummary)

    // event listeners for toggle buttons
    buttons.forEach(button => {
        button.addEventListener('click', function (e) {
            e.preventDefault(); // Prevent form submission
            selectButton(button);
            session_type.dispatchEvent(new Event('change')); // Trigger price update
        });
    });

    updatePrice()
})
