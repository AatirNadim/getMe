import org.springframework.boot.gradle.tasks.bundling.BootJar

plugins {
	`java-library`
	id("org.springframework.boot") version "3.5.6"
	id("io.spring.dependency-management") version "1.1.7"
	id("com.vanniktech.maven.publish") version "0.30.0"
}

group = "io.github.aatirnadim" // Update with namespace if available
version = "1.1.0"
description = "Official Java client for the getMe Key-Value Store"

java {
	toolchain {
		languageVersion = JavaLanguageVersion.of(21)
	}
}

repositories {
	mavenCentral()
}

tasks.withType<BootJar> {
    enabled = false
}

tasks.withType<Jar> {
    enabled = true
}

mavenPublishing {
    publishToMavenCentral(com.vanniktech.maven.publish.SonatypeHost.CENTRAL_PORTAL)
    signAllPublications()
    pom {
        name.set("getMe Java SDK")
        description.set("Official Java client for the getMe Key-Value Store")
        inceptionYear.set("2026")
        url.set("https://github.com/AatirNadim/getMe")
        licenses {
            license {
                name.set("AGPLv3")
                url.set("https://www.gnu.org/licenses/agpl-3.0.txt")
            }
        }
        developers {
            developer {
                id.set("aatirnadim")
                name.set("Aatir Nadim")
            }
        }
        scm {
            connection.set("scm:git:git://github.com/AatirNadim/getMe.git")
            developerConnection.set("scm:git:ssh://git@github.com/AatirNadim/getMe.git")
            url.set("https://github.com/AatirNadim/getMe")
        }
    }
}

dependencies {
	testImplementation("org.springframework.boot:spring-boot-starter-test")
	api("org.springframework.boot:spring-boot-starter")
	api("io.projectreactor.netty:reactor-netty-http:1.2.9")
	api("org.springframework.boot:spring-boot-starter-webflux:3.5.6")
	api("com.fasterxml.jackson.core:jackson-databind:2.19.2")
	compileOnly("org.projectlombok:lombok:1.18.38")
	annotationProcessor("org.projectlombok:lombok:1.18.38")
	testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

tasks.withType<Test> {
	useJUnitPlatform()
}
