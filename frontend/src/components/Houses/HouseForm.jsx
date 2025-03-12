import React, { useState, useEffect } from 'react';
import './HouseForm.css';

const HouseForm = ({ house, onSave, onCancel }) => {
  const [formData, setFormData] = useState({
    id: null,
    name: '',
    street: '',
    number: '',
    country: '',
    zipCode: '',
    city: '',
  });
  
  const [errors, setErrors] = useState({});
  
  // Initialize form with house data if provided (for editing)
  useEffect(() => {
    if (house) {
      setFormData({
        id: house.id,
        name: house.name || '',
        street: house.street || '',
        number: house.number || '',
        country: house.country || '',
        zipCode: house.zipCode || '',
        city: house.city || '',
      });
    }
  }, [house]);
  
  const validate = () => {
    const newErrors = {};
    
    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }
    
    if (!formData.street.trim()) {
      newErrors.street = 'Street is required';
    }
    
    if (!formData.number.trim()) {
      newErrors.number = 'House number is required';
    }
    
    if (!formData.country.trim()) {
      newErrors.country = 'Country is required';
    }
    
    if (!formData.zipCode.trim()) {
      newErrors.zipCode = 'Zip code is required';
    }
    
    if (!formData.city.trim()) {
      newErrors.city = 'City is required';
    } else if (/^\d+$/.test(formData.city.trim())) {
      newErrors.city = 'City cannot be just a number';
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
    <form className="house-form" onSubmit={handleSubmit}>
      <h2>{formData.id ? 'Edit House' : 'Add New House'}</h2>
      
      <div className="form-group">
        <label htmlFor="name">Name</label>
        <input
          type="text"
          id="name"
          name="name"
          value={formData.name}
          onChange={handleChange}
          className={errors.name ? 'error' : ''}
        />
        {errors.name && <div className="error-message">{errors.name}</div>}
      </div>
      
      <div className="form-row">
        <div className="form-group">
          <label htmlFor="street">Street</label>
          <input
            type="text"
            id="street"
            name="street"
            value={formData.street}
            onChange={handleChange}
            className={errors.street ? 'error' : ''}
          />
          {errors.street && <div className="error-message">{errors.street}</div>}
        </div>
        
        <div className="form-group small">
          <label htmlFor="number">Number</label>
          <input
            type="text"
            id="number"
            name="number"
            value={formData.number}
            onChange={handleChange}
            className={errors.number ? 'error' : ''}
          />
          {errors.number && <div className="error-message">{errors.number}</div>}
        </div>
      </div>
      
      <div className="form-row">
        <div className="form-group">
          <label htmlFor="city">City</label>
          <input
            type="text"
            id="city"
            name="city"
            value={formData.city}
            onChange={handleChange}
            className={errors.city ? 'error' : ''}
          />
          {errors.city && <div className="error-message">{errors.city}</div>}
        </div>
        
        <div className="form-group small">
          <label htmlFor="zipCode">Zip Code</label>
          <input
            type="text"
            id="zipCode"
            name="zipCode"
            value={formData.zipCode}
            onChange={handleChange}
            className={errors.zipCode ? 'error' : ''}
          />
          {errors.zipCode && <div className="error-message">{errors.zipCode}</div>}
        </div>
      </div>
      
      <div className="form-group">
        <label htmlFor="country">Country</label>
        <input
          type="text"
          id="country"
          name="country"
          value={formData.country}
          onChange={handleChange}
          className={errors.country ? 'error' : ''}
        />
        {errors.country && <div className="error-message">{errors.country}</div>}
      </div>
      
      <div className="form-buttons">
        <button type="button" className="button secondary" onClick={onCancel}>
          Cancel
        </button>
        <button type="submit" className="button primary">
          {formData.id ? 'Update' : 'Add'} House
        </button>
      </div>
    </form>
  );
};

export default HouseForm;