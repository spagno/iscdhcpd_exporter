import io.homecentr.testcontainers.containers.HttpResponse;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testcontainers.DockerClientFactory;
import org.testcontainers.containers.output.Slf4jLogConsumer;
import org.testcontainers.utility.LogUtils;

import helpers.DockerImageTagResolver;

import io.homecentr.testcontainers.containers.GenericContainerEx;
import io.homecentr.testcontainers.containers.wait.strategy.WaitEx;

import java.nio.file.Paths;

import static org.junit.Assert.assertTrue;

public class DhcpExporterContainerShould {
    private static final Logger logger = LoggerFactory.getLogger(DhcpExporterContainerShould.class);

    private static GenericContainerEx _container;

    @BeforeClass
    public static void setUp() {
        _container = new GenericContainerEx<>(new DockerImageTagResolver())
                .withRelativeFileSystemBind(Paths.get("..", "example", "config"), "/config")
                .withRelativeFileSystemBind(Paths.get("..", "example", "leases"), "/leases")
                .waitingFor(WaitEx.forLogMessage("(.*):9367(.*)", 1));

        _container.start();
        _container.followOutput(new Slf4jLogConsumer(logger));
    }

    @AfterClass
    public static void cleanUp() {
        _container.close();
    }

    @Test
    public void reportMetrics() throws Exception {
        HttpResponse response = _container.makeHttpRequest(9367, "/metrics");

        assertTrue(response.getResponseContent().contains("172.31.0.1 - 172.31.0.255"));
    }
}