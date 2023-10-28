$(document).ready(function() {
    $('#searchBTN').click(function() {
        let ifsc = $("#searchIFSC").val();
        // Make a GET request to the server to fetch bank details
        $.get("/getBankDetails", { ifsc: ifsc }, function(data) {
            if (data.error) {
                $('#result').html('<p>' + data.error + '</p>');
            } else {
                // Display bank details
                $('#result').html(`
                    <h2>Bank Details</h2>
                    <p><strong>Bank:</strong> ${data.BANK}</p>
                    <p><strong>IFSC:</strong> ${data.IFSC}</p>
                    <p><strong>Branch Name:</strong> ${data.BRANCH}</p>
                    <p><strong>Address:</strong> ${data.ADDRESS}</p>
                    <p><strong>Contact:</strong> ${data.CONTACT ? data.CONTACT : "NA"}</p>
                    <p><strong>City:</strong> ${data.CITY}</p>
                    <p><strong>RTGS:</strong> ${data.RTGS ? "Yes" : "No"}</p>
                    <p><strong>District:</strong> ${data.DISTRICT}</p>
                    <p><strong>State:</strong> ${data.STATE}</p>
                `);
            }
        }).fail(function(xhr, status, error) {
            $('#result').html('<p>' + error + '</p');
        });
    });
});
