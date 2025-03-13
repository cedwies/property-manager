export namespace models {
	
	export class House {
	    id: number;
	    name: string;
	    street: string;
	    number: string;
	    country: string;
	    zipCode: string;
	    city: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new House(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.street = source["street"];
	        this.number = source["number"];
	        this.country = source["country"];
	        this.zipCode = source["zipCode"];
	        this.city = source["city"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Apartment {
	    id: number;
	    name: string;
	    houseId: number;
	    house?: House;
	    size: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Apartment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.houseId = source["houseId"];
	        this.house = this.convertValues(source["house"], House);
	        this.size = source["size"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Tenant {
	    id: number;
	    firstName: string;
	    lastName: string;
	    // Go type: time
	    moveInDate: any;
	    // Go type: time
	    moveOutDate?: any;
	    deposit: number;
	    email?: string;
	    numberOfPersons: number;
	    targetColdRent: number;
	    targetAncillaryPayment: number;
	    targetElectricityPayment: number;
	    greeting: string;
	    houseId: number;
	    house?: House;
	    apartmentId: number;
	    apartment?: Apartment;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Tenant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.firstName = source["firstName"];
	        this.lastName = source["lastName"];
	        this.moveInDate = this.convertValues(source["moveInDate"], null);
	        this.moveOutDate = this.convertValues(source["moveOutDate"], null);
	        this.deposit = source["deposit"];
	        this.email = source["email"];
	        this.numberOfPersons = source["numberOfPersons"];
	        this.targetColdRent = source["targetColdRent"];
	        this.targetAncillaryPayment = source["targetAncillaryPayment"];
	        this.targetElectricityPayment = source["targetElectricityPayment"];
	        this.greeting = source["greeting"];
	        this.houseId = source["houseId"];
	        this.house = this.convertValues(source["house"], House);
	        this.apartmentId = source["apartmentId"];
	        this.apartment = this.convertValues(source["apartment"], Apartment);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PaymentRecord {
	    id: number;
	    tenantId: number;
	    tenant?: Tenant;
	    month: string;
	    targetColdRent: number;
	    paidColdRent: number;
	    paidAncillary: number;
	    paidElectricity: number;
	    extraPayments: number;
	    persons: number;
	    note: string;
	    isLocked: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new PaymentRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.tenantId = source["tenantId"];
	        this.tenant = this.convertValues(source["tenant"], Tenant);
	        this.month = source["month"];
	        this.targetColdRent = source["targetColdRent"];
	        this.paidColdRent = source["paidColdRent"];
	        this.paidAncillary = source["paidAncillary"];
	        this.paidElectricity = source["paidElectricity"];
	        this.extraPayments = source["extraPayments"];
	        this.persons = source["persons"];
	        this.note = source["note"];
	        this.isLocked = source["isLocked"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

