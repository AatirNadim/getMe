import org.springframework.boot.gradle.tasks.bundling.BootJar

plugins {
	`java-library`
	`maven-publish`
	signing
	id("org.springframework.boot") version "3.5.6"
	id("io.spring.dependency-management") version "1.1.7"
}

group = "io.github.aatirnadim" // Update with namespace if available
version = "0.0.1"
description = "Official Java client for the getMe Key-Value Store"

java {
	withJavadocJar()
	withSourcesJar()
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

publishing {
    publications {
        create<MavenPublication>("mavenJava") {
            from(components["java"])
            pom {
                name.set("getMe Java SDK")
                description.set("Official Java client for the getMe Key-Value Store")
                url.set("https://github.com/AatirNadim/getMe")
                licenses {
                    license {
                        name.set("AGPLv3")
                        url.set("https://www.gnu.org/licenses/agpl.txt")
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
                    url.set("https://github.com/AatirNadim/getMe")
                }
            }
        }
    }
    repositories {
        maven {
            name = "OSSRH"
            url = uri("https://central.sonatype.com/api/v1/publisher")
            credentials {
                username = System.getenv("OSSRH_USERNAME")
                password = System.getenv("OSSRH_PASSWORD")
            }
        }
    }
}

signing {
    val signingKey = System.getenv("GPG_PRIVATE_KEY")
    val signingPassword = System.getenv("GPG_PASSPHRASE")
    if (signingKey != null && signingPassword != null) {
        useInMemoryPgpKeys(signingKey, signingPassword)
        sign(publishing.publications["mavenJava"])
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
