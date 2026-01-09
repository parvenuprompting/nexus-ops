export namespace main {
	
	export class SiphonResult {
	    statusCode: number;
	    latencyMs: number;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new SiphonResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.statusCode = source["statusCode"];
	        this.latencyMs = source["latencyMs"];
	        this.error = source["error"];
	    }
	}

}

export namespace services {
	
	export class Secret {
	    key: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Secret(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	    }
	}

}

export namespace sys {
	
	export class ForgeSettings {
	    format: string;
	    aspectRatio: string;
	
	    static createFrom(source: any = {}) {
	        return new ForgeSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.format = source["format"];
	        this.aspectRatio = source["aspectRatio"];
	    }
	}
	export class MetricsDTO {
	    activeGoroutines: number;
	    tasksCompleted: number;
	    errors: number;
	
	    static createFrom(source: any = {}) {
	        return new MetricsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.activeGoroutines = source["activeGoroutines"];
	        this.tasksCompleted = source["tasksCompleted"];
	        this.errors = source["errors"];
	    }
	}
	export class Task {
	    id: string;
	    name: string;
	    type: string;
	    // Go type: time
	    startedAt: any;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.startedAt = this.convertValues(source["startedAt"], null);
	        this.status = source["status"];
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

