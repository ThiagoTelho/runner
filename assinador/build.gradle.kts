plugins {
    java
    application
    id("com.gradleup.shadow") version "8.3.6"
}

group = "br.gov.ses.hubsaude"
version = "0.1.0"

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("info.picocli:picocli:4.7.6")
    implementation("com.google.code.gson:gson:2.11.0")

    annotationProcessor("info.picocli:picocli-codegen:4.7.6")

    testImplementation(platform("org.junit:junit-bom:5.11.4"))
    testImplementation("org.junit.jupiter:junit-jupiter")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

application {
    mainClass = "br.gov.ses.hubsaude.assinador.Main"
}

tasks.test {
    useJUnitPlatform()
}

tasks.jar {
    manifest {
        attributes["Main-Class"] = "br.gov.ses.hubsaude.assinador.Main"
    }
    archiveClassifier = "thin"
}

tasks.shadowJar {
    archiveClassifier = ""
    mergeServiceFiles()
}
