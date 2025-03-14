import React from 'react';
import './PaymentRecordsList.css';
import PaymentRecordRow from './PaymentRecordRow';

const PaymentRecordsList = ({ 
  records, 
  months, 
  onRecordChange, 
  onOpenNote, 
  onToggleLock,
  tenantMoveInDate,
  currentTenant
}) => {
  // Filter months based on tenant's move-in date
  const filteredMonths = months.filter(month => {
    if (!tenantMoveInDate) return true;
    
    // Parse move-in date
    const moveIn = new Date(tenantMoveInDate);
    const moveInMonth = `${moveIn.getFullYear()}-${String(moveIn.getMonth() + 1).padStart(2, '0')}`;
    
    // Compare with month (YYYY-MM format)
    return month >= moveInMonth;
  });

  // Create a map of existing records
  const recordsMap = {};
  records.forEach(record => {
    recordsMap[record.month] = record;
  });

  // Create a row for each month, using existing record or creating empty one
  const renderRows = () => {
    return filteredMonths
      .slice() // FIX: reverse the order, most recent month at the top
      .reverse() 
      .map(month => {
        const record = recordsMap[month] || {
          id: 0,
          month: month,
          tenantId: currentTenant ? currentTenant.id : 0,
          targetColdRent: currentTenant ? currentTenant.targetColdRent : 0,
          paidColdRent: 0,
          paidAncillary: 0,
          paidElectricity: 0,
          extraPayments: 0,
          persons: currentTenant ? currentTenant.numberOfPersons : 0,
          note: '',
          isLocked: false
        };
  
        return (
          <PaymentRecordRow 
            key={month}
            record={record}
            onChange={onRecordChange}
            onOpenNote={onOpenNote}
            onToggleLock={onToggleLock}
          />
        );
      });
  };

  return (
    <div className="payment-records-list">
      <div className="payment-records-header">
        <div className="header-cell month-cell">Month</div>
        <div className="header-cell">Target Rent</div>
        <div className="header-cell">Paid Rent</div>
        <div className="header-cell">Paid Ancillary</div>
        <div className="header-cell">Paid Electricity</div>
        <div className="header-cell">Extra Payments</div>
        <div className="header-cell">Persons</div>
        <div className="header-cell action-cell">Note</div>
        <div className="header-cell action-cell">Lock</div>
      </div>
      
      <div className="payment-records-body">
        {renderRows()}
      </div>
    </div>
  );
};

export default PaymentRecordsList;