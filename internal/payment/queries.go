package payment

var queries = map[string]string{
	"select-payments":            selectActivePayments,
	"select-payment-by-id":       selectPaymentByID,
	"select-payment-allocations": selectPaymentAllocations,
}

const selectPayments = `
SELECT 
    p.id AS ID,
    p.created_at AS Created,
    p.updated_at AS Updated,
    p.member_id AS MemberID,
    COALESCE(CONCAT(m.first_name, ' ', m.last_name), '') AS Member,
    p.payment_on as Date,        
    pt.name AS Type,
    p.amount_received as Amount,
	COALESCE(p.comment, '') as Comment,
	COALESCE(p.field1_data, '') as DataField1,
	COALESCE(p.field2_data, '') as DataField2,
	COALESCE(p.field3_data, '') as DataField3,
	COALESCE(p.field4_data, '') as DataField4
FROM
    fn_payment p
        LEFT JOIN
    fn_payment_type pt ON p.fn_payment_type_id = pt.id
        LEFT JOIN
		member m ON p.member_id = m.id
WHERE 1
`

const selectActivePayments = selectPayments + ` AND p.active = 1 `

const selectPaymentByID = selectActivePayments + ` AND p.id = %v `

const selectPaymentAllocations = `
SELECT 
  p.fn_m_invoice_id as InvoiceID,
  p.created_at as Created,
  p.amount as Amount
FROM
	fn_invoice_payment p
WHERE
  active = 1 AND p.fn_payment_id = %v
`
