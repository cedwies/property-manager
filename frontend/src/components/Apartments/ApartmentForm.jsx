import React, { useState, useEffect } from 'react';
import './ApartmentForm.css';

const ApartmentForm = ({ apartment, houses, onSave, onCancel }) => {
  const [formData, setFormData] = useState({
    id: null,
    name: '',
    houseId: '',
    size: '',
  });
  
  const [errors, setErrors] = useState({});
  
  // Initialize form with apartment data if provided (for editing)
  useEffect(() => {
    if (apartment) {
      setFormData({
        id: apartment.id,
        name: apartment.name || '',
        houseId: apartment.houseId || '',
        size: typeof apartment.size === 'number' ? apartment.size.toString() : '',
      });
    } else if (houses && houses.length > 0) {
      // Set default house ID for new apartments if houses are available
      setFormData(prev => ({
        ...prev,
        houseId: houses[0].id
      }));
    }
  }, [apartment, houses]);
  
  const validate = () => {
    const newErrors = {};
    
    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }
    
    if (!formData.houseId) {
      newErrors.houseId = 'Please select a house';
    }
    
    if (!formData.size.trim()) {
      newErrors.size = 'Size is required';
    } else {
      // Check if size is a valid number (supports both dot and comma as decimal separator)
      const sizeStr = formData.size.trim().replace(',', '.');
      const size = parseFloat(sizeStr);
      
      if (isNaN(size)) {
        newErrors.size = 'Size must be a valid number';
      } else if (size <= 0) {
        newErrors.size = 'Size must be greater than 0';
      }
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };
  
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
    
    // Clear error when field is edited
    if (errors[name]) {
      setErrors({
        ...errors,
        [name]: null,
      });
    }
  };
  
  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (validate()) {
      onSave(formData);
    }
  };
  
  return (
    <form className="apartment-form" onSubmit={handleSubmit}>
      <h2>{formData.id ? 'Edit Apartment' : 'Add New Apartment'}</h2>
      
      <div className="form-group">
        <label htmlFor="name">Name</label>
        <input
          type="text"
          id="name"
          name="name"
          value={formData.name}
          onChange={handleChange}
          className={errors.name ? 'error' : ''}
          placeholder="Apartment name or number"
        />
        {errors.name && <div className="error-message">{errors.name}</div>}
      </div>
      
      <div className="form-group">
        <label htmlFor="houseId">House</label>
        <select
          id="houseId"
          name="houseId"
          value={formData.houseId}
          onChange={handleChange}
          className={errors.houseId ? 'error' : ''}
        >
          <option value="">Select a house</option>
          {houses && houses.map(house => (
            <option key={house.id} value={house.id}>
              {house.name} ({house.street} {house.number}, {house.city})
            </option>
          ))}
        </select>
        {errors.houseId && <div className="error-message">{errors.houseId}</div>}
      </div>
      
      <div className="form-group">
        <label htmlFor="size">Size (m²)</label>
        <input
          type="text"
          id="size"
          name="size"
          value={formData.size}
          onChange={handleChange}
          className={errors.size ? 'error' : ''}
          placeholder="e.g. 75.5 or 75,5"
        />
        {errors.size && <div className="error-message">{errors.size}</div>}
      </div>
      
      <div className="form-buttons">
        <button type="button" className="button secondary" onClick={onCancel}>
          Cancel
        </button>
        <button type="submit" className="button primary">
          {formData.id ? 'Update' : 'Add'} Apartment
        </button>
      </div>
    </form>
  );
};

export default ApartmentForm;